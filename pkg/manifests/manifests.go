package manifests

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rancher/k3os/pkg/config"
	"github.com/sirupsen/logrus"
)

const (
	retryMax     = 5
	manifestsDir = "/var/lib/rancher/k3s/server/manifests"
)

func ApplyBootManifests(cfg *config.CloudConfig) error {
	manifests := cfg.BootManifests
	if len(manifests) == 0 {
		return nil
	}
	filesToWrite := make(map[string][]byte)
	var err error
	for _, m := range manifests {
		var data []byte
		retries := 0
		for retryMax > retries {
			resp, err := http.Get(m.URL)
			if err != nil {
				logrus.Errorf("manifest download failed for %q, retrying [%d/%d]", m.URL, retries, retryMax)
				retries++
				continue
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				logrus.Errorf("manifest download returned non-200 status code for %q, retrying [%d/%d]", m.URL, retries, retryMax)
				retries++
				continue
			}
			data, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				errors.Wrap(err, "failed reading manifest URL body")
				return err
			}
			if len(data) == 0 {
				return errors.New("empty manifest for %q")
			}
			if m.SHA256 != "" {
				sum := sha256.Sum256(data)
				if fmt.Sprintf("%x", sum) != m.SHA256 {
					return fmt.Errorf("sha256 failed for manifest: %s", m.URL)
				}
			}
			name := m.URL[strings.LastIndex(m.URL, "/")+1:]
			filesToWrite[name] = data
			break
		}
	}
	for file, data := range filesToWrite {
		p := filepath.Join(manifestsDir, file)
		if err := ioutil.WriteFile(p, data, 0600); err != nil {
			return err
		}

	}
	return err
}
