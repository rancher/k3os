package control

import (
	"github.com/niusmallnan/k3os/pkg/util"
	"os"

	"github.com/niusmallnan/k3os/config"
	"github.com/niusmallnan/k3os/pkg/command"
	pkgHostname "github.com/niusmallnan/k3os/pkg/hostname"
	"github.com/niusmallnan/k3os/pkg/module"
	"github.com/niusmallnan/k3os/pkg/ssh"
	"github.com/niusmallnan/k3os/pkg/sysctl"

	"github.com/sirupsen/logrus"
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
	// setup hostname
	if err := pkgHostname.SetHostname(cfg); err != nil {
		return err
	}
	// setup /etc/hosts
	if err := pkgHostname.SyncHostname(); err != nil {
		return err
	}
	// setup ssh authorized_keys
	for _, username := range config.SSHUsers {
		if err := ssh.SetAuthorizedKeys(username, cfg); err != nil {
			logrus.Error(err)
		}
	}
	// setup kernel modules
	if err := module.LoadModules(cfg); err != nil {
		return err
	}
	// setup sysctl
	if err := sysctl.ConfigureSysctl(cfg); err != nil {
		return err
	}
	// run command
	if err := command.ExecuteCommand(cfg.Runcmd); err != nil {
		return err
	}
	// run rc.local
	return util.RunScript("/etc/rc.local")
}
