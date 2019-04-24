package main

import (
	"fmt"
	"os"

	"github.com/rancher/k3os/pkg/cliinstall"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	debug = true
)

func main() {
	app := cli.NewApp()

	app.Name = "k3os config"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Destination: &debug,
			EnvVar:      "K3OS_DEBUG",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(_ *cli.Context) error {
	if os.Getuid() != 0 {
		return fmt.Errorf("must run %s as root", os.Args[0])
	}
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	return cliinstall.Run()
}
