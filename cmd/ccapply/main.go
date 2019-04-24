package main

import (
	"fmt"
	"os"

	"github.com/rancher/k3os/pkg/cliinstall"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	initrd = false
	boot   = false
	config = false
)

func main() {
	app := cli.NewApp()

	app.Name = "k3os config"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "init",
			Destination: &initrd,
			Usage:       "Run initrd stage",
		},
		cli.BoolFlag{
			Name:        "boot",
			Destination: &boot,
			Usage:       "Run boot stage",
		},
		cli.BoolFlag{
			Name:        "config",
			Destination: &config,
			Usage:       "Run os-config stage",
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
	return cliinstall.Run()
}
