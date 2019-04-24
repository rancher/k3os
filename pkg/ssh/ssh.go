package ssh

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/rancher/k3os/pkg/config"
	"github.com/rancher/k3os/pkg/util"
	"github.com/sirupsen/logrus"
)

const (
	sshDir         = ".ssh"
	authorizedFile = "authorized_keys"
)

func SetAuthorizedKeys(cfg *config.CloudConfig, withNet bool) error {
	bytes, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		return err
	}
	uid, gid, homeDir, err := findUserHomeDir(bytes, "rancher")
	if err != nil {
		return err
	}
	userSSHDir := path.Join(homeDir, sshDir)
	if _, err := os.Stat(userSSHDir); os.IsNotExist(err) {
		if err = os.Mkdir(userSSHDir, 0700); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	if err = os.Chown(userSSHDir, uid, gid); err != nil {
		return err
	}
	userAuthorizedFile := path.Join(userSSHDir, authorizedFile)
	for _, key := range cfg.SSHAuthorizedKeys {
		if err = authorizeSSHKey(key, userAuthorizedFile, uid, gid, withNet); err != nil {
			logrus.Errorf("failed to authorize SSH key %s: %v", key, err)
		}
	}
	return nil
}

func getKey(key string, withNet bool) (string, error) {
	if !strings.HasPrefix(key, "github:") {
		return key, nil
	}

	if !withNet {
		return "", nil
	}

	resp, err := http.Get(fmt.Sprintf("https://github.com/%s.keys", strings.TrimPrefix(key, "github:")))
	if err != nil {
		return "", err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	return string(bytes), err
}

func authorizeSSHKey(key, file string, uid, gid int, withNet bool) error {
	key, err := getKey(key, withNet)
	if err != nil || key == "" {
		return err
	}

	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		if err = f.Chmod(0600); err != nil {
			return err
		}
		if err = f.Close(); err != nil {
			return err
		}
		info, err = os.Stat(file)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if !strings.Contains(string(bytes), key) {
		bytes = append(bytes, []byte(key)...)
		bytes = append(bytes, '\n')
	}
	perm := info.Mode().Perm()
	if err = util.WriteFileAtomic(file, bytes, perm); err != nil {
		return err
	}
	return os.Chown(file, uid, gid)
}

func findUserHomeDir(bytes []byte, username string) (uid, gid int, homeDir string, err error) {
	for _, line := range strings.Split(string(bytes), "\n") {
		if strings.HasPrefix(line, username) {
			split := strings.Split(line, ":")
			if len(split) < 6 {
				break
			}
			uid, err = strconv.Atoi(split[2])
			if err != nil {
				return -1, -1, "", err
			}
			gid, err = strconv.Atoi(split[3])
			if err != nil {
				return -1, -1, "", err
			}
			homeDir = split[5]
		}
	}
	return
}
