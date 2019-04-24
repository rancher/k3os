package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func Convert(from, to interface{}) error {
	bytes, err := yaml.Marshal(from)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bytes, to)
}

func Env2Map(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, s := range env {
		d := strings.Split(s, "=")
		m[d[0]] = d[1]
	}
	return m
}

func MapCopy(data map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for k, v := range data {
		result[k] = Copy(v)
	}
	return result
}

func Copy(d interface{}) interface{} {
	switch d := d.(type) {
	case map[string]interface{}:
		return MapCopy(d)
	case []interface{}:
		return SliceCopy(d)
	default:
		return d
	}
}

func SliceCopy(data []interface{}) []interface{} {
	result := make([]interface{}, len(data), len(data))
	for k, v := range data {
		result[k] = Copy(v)
	}
	return result
}

func FilterKeys(data map[string]interface{}, key []string) (filtered, rest map[string]interface{}) {
	if len(key) == 0 {
		return data, map[string]interface{}{}
	}
	filtered = map[string]interface{}{}
	rest = MapCopy(data)
	k := key[0]
	if d, ok := data[k]; ok {
		switch d := d.(type) {
		case map[string]interface{}:
			f, r := FilterKeys(d, key[1:])
			if len(f) != 0 {
				filtered[k] = f
			}
			if len(r) != 0 {
				rest[k] = r
			} else {
				delete(rest, k)
			}
		default:
			filtered[k] = d
			delete(rest, k)
		}
	}
	return
}

func WriteFileAtomic(filename string, data []byte, perm os.FileMode) error {
	dir, file := path.Split(filename)
	tempFile, err := ioutil.TempFile(dir, fmt.Sprintf(".%s", file))
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())
	if _, err := tempFile.Write(data); err != nil {
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tempFile.Name(), perm); err != nil {
		return err
	}
	return os.Rename(tempFile.Name(), filename)
}

func WriteToFile(data interface{}, filename string) error {
	content, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(filename), os.ModeDir|0700); err != nil {
		return err
	}
	return WriteFileAtomic(filename, content, 400)
}

func HTTPDownloadToFile(url, dest string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return WriteFileAtomic(dest, body, 0644)
}

func HTTPLoadBytes(url string) ([]byte, error) {
	var resp *http.Response
	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("non-200 http response: %d", resp.StatusCode)
		}

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return bytes, nil
	}

	return nil, err
}

func ExistsAndExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	mode := info.Mode().Perm()
	return mode&os.ModePerm != 0
}

func RunScript(path string, arg ...string) error {
	if !ExistsAndExecutable(path) {
		return nil
	}

	script, err := os.Open(path)
	if err != nil {
		return err
	}

	magic := make([]byte, 2)
	if _, err = script.Read(magic); err != nil {
		return err
	}

	cmd := exec.Command("/bin/sh", path)
	if string(magic) == "#!" {
		cmd = exec.Command(path, arg...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func FileCopy(src, dest string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return WriteFileAtomic(dest, data, 0644)
}

func UnescapeKernelParams(s string) string {
	s = strings.Replace(s, `\"`, `"`, -1)
	s = strings.Replace(s, `\'`, `'`, -1)
	return s
}

func EnsureDirectoryExists(dir string) error {
	info, err := os.Stat(dir)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("%s is not a directory", dir)
		}
	} else {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
