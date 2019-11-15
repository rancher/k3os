package upgrade

import (
	"encoding/json"
	"io"

	"github.com/rancher/k3os/pkg/util"
)

type github struct {
	repo string
}

func (g *github) Latest() (rel Release, err error) {
	err = util.Stream(`https://api.github.com/repos/`+g.repo+`/releases/latest`, func(body io.Reader) error {
		if err := json.NewDecoder(body).Decode(&rel); err != nil {
			return err
		}
		return nil
	})
	return rel, err
}

func (g *github) Release(tag string) (rel Release, err error) {
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
