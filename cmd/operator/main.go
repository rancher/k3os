package main

import (
	"os"

	"github.com/rancher/k3os/cmd/operator/agent"
	"github.com/rancher/k3os/cmd/operator/upgrade"
	"github.com/rancher/k3os/pkg/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "k3os-operator"
	app.Version = version.Version
	app.Commands = []cli.Command{
		agent.Command,
		upgrade.Command,
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			EnvVar: "K3OS_DEBUG",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.SetReportCaller(true)
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
