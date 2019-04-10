package netinit

import (
	"os"

	"github.com/niusmallnan/k3os/config"
	"github.com/niusmallnan/k3os/pkg/network"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func Main() {
	app := cli.NewApp()
	app.Author = "Rancher Labs, Inc."
	app.EnableBashCompletion = true
	app.HideHelp = true
	app.Name = os.Args[0]
	app.Usage = "k3os network init"
	app.Version = config.OSVersion
	app.Action = netInit
	app.Run(os.Args)
}

func netInit(c *cli.Context) error {
	cfg := config.LoadConfig("", false)
	// setup network
	if err := network.SettingNetwork(cfg); err != nil {
		logrus.Fatalf("failed to setting network: %v", err)
	}
	// setup dns
	if err := network.SettingDNS(cfg); err != nil {
		logrus.Fatalf("failed to setting dns: %v", err)
	}
	// setup proxy
	if err := network.SettingProxy(cfg); err != nil {
		logrus.Fatalf("failed to setting proxy: %v", err)
	}
	return nil
}
