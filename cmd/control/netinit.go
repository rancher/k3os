package control

import (
	"os"

	"github.com/niusmallnan/k3os/config"
	"github.com/niusmallnan/k3os/pkg/network"

	"github.com/urfave/cli"
)

func NetInitMain() {
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
	// setup dns
	return network.SettingDNS(cfg)
}
