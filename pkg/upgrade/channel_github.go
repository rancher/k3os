package upgrade

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rancher/k3os/pkg/system"
	"github.com/rancher/k3os/pkg/util"
	"github.com/sirupsen/logrus"
)

type githubChannel struct {
	uri    string
	repo   string
	latest string
}

func (g *githubChannel) resolve() error {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(http.MethodGet, g.uri, nil)
	if err != nil {
		return err
	}

	if uuid, err := ioutil.ReadFile("/sys/class/dmi/id/product_uuid"); err != nil {
		logrus.Warn(err)
	} else {
		req.Header.Set("x-k3os-uuid", string(uuid))
	}
	ver, err := system.GetVersion()
	if err != nil {
		logrus.Warn(err)
	}
	if ver.Runtime != "" {
		req.Header.Set("x-k3os-version", ver.Runtime)
	}
	ver, err = system.GetKernelVersion()
	if err != nil {
		logrus.Warn(err)
	}
	if ver.Runtime != "" {
		req.Header.Set("x-k3os-kernel", ver.Runtime)
	}
	req.Header.Set("x-k3os-arch", runtime.GOARCH)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusFound {
		return fmt.Errorf("expecting redirect: %s %s", res.Proto, res.Status)
	}
	found, err := res.Location()
	if err != nil {
		return err
	}
	if !strings.HasSuffix(found.Hostname(), "github.com") {
		return fmt.Errorf("not a github url: %v", found)
	}
	g.latest = filepath.Base(found.Path)

	p := strings.SplitN(strings.TrimLeft(found.Path, `/`), `/`, 3)
	if len(p) != 3 {
		return fmt.Errorf("unable to parse repo: %v", found.Path)
	}
	g.repo = path.Join(p[0], p[1])
	return nil
}

func (g *githubChannel) Latest() (rel Release, err error) {
	if g.latest == "" {
		if err = g.resolve(); err != nil {
			return rel, err
		}
	}
	if g.latest == "" {
		return rel, fmt.Errorf("unable to resolve latest")
	}
	return g.Release(g.latest)
}

func (g *githubChannel) Release(tag string) (rel Release, err error) {
	if tag == "" || tag == "latest" {
		return g.Latest()
	}
	return rel, util.Stream(`https://api.github.com/repos/`+g.repo+`/releases/tags/`+tag, func(body io.Reader) error {
		if err := json.NewDecoder(body).Decode(&rel); err != nil {
			return err
		}
		return nil
	})
}
