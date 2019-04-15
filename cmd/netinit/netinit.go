package netinit

import (
	"fmt"
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
	app.Usage = fmt.Sprintf("%s K3OS(%s)", app.Name, config.OSBuildDate)
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
	// setup proxy environments
	network.SettingProxyEnvironments(cfg)
	// setup proxy
	if err := network.SettingProxy(); err != nil {
		logrus.Fatalf("failed to setting proxy: %v", err)
	}
	return nil
}
