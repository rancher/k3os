package install

import (
	"fmt"
	"os"

	"github.com/rancher/k3os/pkg/cliinstall"
	"github.com/rancher/k3os/pkg/mode"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func Command() cli.Command {
	mode, _ := mode.Get()
	return cli.Command{
		Name:  "install",
		Usage: "install k3OS",
		Flags: []cli.Flag{},
		Before: func(c *cli.Context) error {
			if os.Getuid() != 0 {
				return fmt.Errorf("must be run as root")
			}
			return nil
		},
		Action: func(*cli.Context) {
			if err := cliinstall.Run(); err != nil {
				logrus.Error(err)
			}
		},
		Hidden: mode == "local",
	}
}
