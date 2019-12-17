package operator

import (
	"github.com/rancher/k3os/pkg/cli/operator/agent"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Command `operator`
func Command() cli.Command {
	return command // value copy on purpose
}

var (
	command = cli.Command{
		Name:      "operator",
		Usage:     "operate k3OS",
		ShortName: "ops",
		Subcommands: []cli.Command{
			agent.Command(),
		},
		Before: func(c *cli.Context) error {
			if c.GlobalBool("debug") {
				logrus.SetLevel(logrus.DebugLevel)
				logrus.SetReportCaller(true)
			}
			return nil
		},
	}
)
