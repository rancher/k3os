package ssh

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/rancher/k3os/config"
	"github.com/rancher/k3os/pkg/util"
	"github.com/sirupsen/logrus"
)

const (
	sshDir         = ".ssh"
	authorizedFile = "authorized_keys"
)

func SetAuthorizedKeys(username string, cfg *config.CloudConfig) error {
	bytes, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		return err
	}
	uid, gid, homeDir, err := findUserHomeDir(bytes, username)
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
	for _, key := range cfg.K3OS.SSH.AuthorizedKeys {
		if err = authorizeSSHKey(key, userAuthorizedFile, uid, gid); err != nil {
			logrus.Errorf("failed to authorize SSH key %s: %v", key, err)
		}
	}
	return nil
}

func SetHostKeys(cfg *config.CloudConfig) error {
	for _, t := range []string{"rsa", "dsa", "ecdsa", "ed25519"} {
		f := fmt.Sprintf("/etc/ssh/ssh_host_%s_key", t)
		p := fmt.Sprintf("/etc/ssh/ssh_host_%s_key.pub", t)
		key, keyExist := cfg.K3OS.SSH.HostKeys[t]
		pub, pubExist := cfg.K3OS.SSH.HostKeys[t+"-pub"]
		if keyExist && pubExist {
			if err := util.WriteFileAtomic(f, []byte(key), 0600); err != nil {
				return err
			}
			if err := util.WriteFileAtomic(p, []byte(pub), 0600); err != nil {
				return err
			}
			continue
		}
		if _, err := os.Stat(f); err != nil || os.IsNotExist(err) {
			continue
		}
		if _, err := os.Stat(p); err != nil || os.IsNotExist(err) {
			continue
		}
		fb, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		pb, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		if err := config.Set(fmt.Sprintf("k3os.ssh.host_keys.%s", t), string(fb)); err != nil {
			return err
		}
		if err := config.Set(fmt.Sprintf("k3os.ssh.host_keys.%s-pub", t), string(pb)); err != nil {
			return err
		}
	}
	return nil
}

func authorizeSSHKey(key, file string, uid, gid int) error {
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
