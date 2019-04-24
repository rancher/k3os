package main

import (
	"os"

	"github.com/rancher/k3os/pkg/cliinstall"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	install bool
)

func main() {
	app := cli.NewApp()

	app.Name = "k3os config"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "install",
			Destination: &install,
		},
	}

	app.Action = run
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(cli *cli.Context) error {
	if install {
		return cliinstall.Run(cli.Args())
	}

	return nil
}
