package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rancher/k3os/pkg/cc"
	"github.com/rancher/k3os/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	initrd      = false
	bootPhase   = false
	configPhase = false
	debug       = false
	dump        = false
	dumpJSON    = false
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
		cli.BoolFlag{
			Name:        "initrd",
			Destination: &initrd,
			Usage:       "Run initrd stage",
		},
		cli.BoolFlag{
			Name:        "boot",
			Destination: &bootPhase,
			Usage:       "Run boot stage",
		},
		cli.BoolFlag{
			Name:        "config",
			Destination: &configPhase,
			Usage:       "Run os-config stage",
		},
		cli.BoolFlag{
			Name:        "dump",
			Destination: &dump,
			Usage:       "Print current configuration",
		},
		cli.BoolFlag{
			Name:        "dump-json",
			Destination: &dumpJSON,
			Usage:       "Print current configuration in json",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(_ *cli.Context) error {
	if err := doRun(); err != nil {
		logrus.Error(err)
	}

	return nil
}

func doRun() error {
	if os.Getuid() != 0 {
		return fmt.Errorf("must run %s as root", os.Args[0])
	}
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	if initrd {
		return cc.InitApply(&cfg)
	} else if bootPhase {
		return cc.BootApply(&cfg)
	} else if configPhase {
		return cc.ConfigApply(&cfg)
	} else if dump {
		return config.Write(cfg, os.Stdout)
	} else if dumpJSON {
		return json.NewEncoder(os.Stdout).Encode(&cfg)
	}

	return cc.RunApply(&cfg)
}
