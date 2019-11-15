package upgrade

import "github.com/rancher/k3os/pkg/upgrade"

type Channel struct {
	upchan upgrade.Channel
	uri    string
}

func (c *Channel) Latest() (upgrade.Release, error) {
	return c.upchan.Latest()
}

func (c *Channel) Release(tag string) (upgrade.Release, error) {
	return c.upchan.Release(tag)
}

func (c *Channel) Set(uri string) (err error) {
	c.uri = uri
	c.upchan, err = upgrade.NewChannel(uri)
	return err
}

func (c *Channel) String() string {
	return c.uri
}
