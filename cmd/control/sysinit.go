package control

import (
	"os"

	"github.com/niusmallnan/k3os/config"
	pkgHostname "github.com/niusmallnan/k3os/pkg/hostname"

	"github.com/urfave/cli"
)

func SysInitMain() {
	app := cli.NewApp()
	app.Author = "Rancher Labs, Inc."
	app.EnableBashCompletion = true
	app.HideHelp = true
	app.Name = os.Args[0]
	app.Usage = "k3os system init"
	app.Version = config.OSVersion
	app.Action = sysInit
	app.Run(os.Args)
}

func sysInit(c *cli.Context) error {
	cfg := config.LoadConfig("", false)
	// setup k3os hostname
	err := pkgHostname.SetHostname(cfg)
	if err != nil {
		return err
	}
	// setup k3os /etc/hosts
	return pkgHostname.SyncHostname()
}
