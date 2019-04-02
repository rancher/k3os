package util

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func Convert(from, to interface{}) error {
	bytes, err := yaml.Marshal(from)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bytes, to)
}

func GenTemplate(in io.Reader, out io.Writer) error {
	bytes, err := ioutil.ReadAll(in)
	if err != nil {
		logrus.Fatal("could not read from stdin")
	}
	tpl := template.Must(template.New("k3osConfig").Parse(string(bytes)))
	return tpl.Execute(out, Env2Map(os.Environ()))
}

func Env2Map(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, s := range env {
		d := strings.Split(s, "=")
		m[d[0]] = d[1]
	}
	return m
}

func MapCopy(data map[interface{}]interface{}) map[interface{}]interface{} {
	result := map[interface{}]interface{}{}
	for k, v := range data {
		result[k] = Copy(v)
	}
	return result
}

func Copy(d interface{}) interface{} {
	switch d := d.(type) {
	case map[interface{}]interface{}:
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

func FilterKeys(data map[interface{}]interface{}, key []string) (filtered, rest map[interface{}]interface{}) {
	if len(key) == 0 {
		return data, map[interface{}]interface{}{}
	}
	filtered = map[interface{}]interface{}{}
	rest = MapCopy(data)
	k := key[0]
	if d, ok := data[k]; ok {
		switch d := d.(type) {
		case map[interface{}]interface{}:
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

func GetValue(args string, data map[interface{}]interface{}) (interface{}, map[interface{}]interface{}) {
	parts := strings.Split(args, ".")
	tempData := data
	for index, part := range parts {
		last := index+1 == len(parts)
		val, ok := tempData[part]
		if !ok {
			break
		}
		if last {
			return val, tempData
		}
		d, ok := val.(map[interface{}]interface{})
		if !ok {
			break
		}
		tempData = d
	}
	return "", tempData
}

func SetValue(args string, data map[interface{}]interface{}, value interface{}) (interface{}, map[interface{}]interface{}) {
	parts := strings.Split(args, ".")
	copyData := MapCopy(data)
	tempData := copyData
	for index, part := range parts {
		last := index+1 == len(parts)
		val, ok := tempData[part]
		if last {
			if v, ok := value.(string); ok {
				value = UnmarshalValue(v)
			}
			tempData[part] = value
			return value, copyData
		}
		if !last && !ok {
			d := map[interface{}]interface{}{}
			tempData[part] = d
			tempData = d
			continue
		}
		if !ok {
			break
		}
		if last {
			return val, copyData
		}
		d, ok := val.(map[interface{}]interface{})
		if !ok {
			break
		}
		tempData = d
	}
	return "", copyData
}

func UnmarshalValue(value string) (result interface{}) {
	if err := yaml.Unmarshal([]byte(value), &result); err != nil {
		result = value
	}
	return
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

func Yes(question string) bool {
	fmt.Printf("%s [y/N]: ", question)
	in := bufio.NewReader(os.Stdin)
	line, err := in.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return strings.ToLower(line[0:1]) == "y"
}
