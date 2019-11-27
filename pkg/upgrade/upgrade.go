package upgrade

import (
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/rancher/k3os/pkg/util"
)

// Channel represents a release stream
type Channel interface {
	Latest() (Release, error)
	Release(name string) (Release, error)
}

// Release is currently a GitHub API release
type Release struct {
	Pre    bool    `json:"prerelease,omitempty"`
	Name   string  `json:"tag_name,omitempty"`
	Assets []Asset `json:"assets,omitempty"`
}

// Asset is currently a GitHub API releases asset
type Asset struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"browser_download_url,omitempty"`
}

// Asset retreives an asset by name
func (rel *Release) Asset(name string) *Asset {
	for _, a := range rel.Assets {
		if a.Name == name {
			return &a
		}
	}
	return nil
}

// KernelVersion returns the kernel version for the Release
func (rel *Release) KernelVersion() (string, error) {
	var (
		kvn = `k3os-kernel-version-` + runtime.GOARCH
		kva = rel.Asset(kvn)
		ver string
	)
	if kva == nil {
		return "", fmt.Errorf("asset not found: %s", kvn)
	}
	return ver, util.Stream(kva.URL, func(body io.Reader) error {
		bytes, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}
		ver = strings.TrimSpace(string(bytes))
		return nil
	})
}
