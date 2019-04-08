package control

import (
	"os"

	"github.com/niusmallnan/k3os/config"
	"github.com/niusmallnan/k3os/pkg/command"
	pkgHostname "github.com/niusmallnan/k3os/pkg/hostname"
	"github.com/niusmallnan/k3os/pkg/module"
	"github.com/niusmallnan/k3os/pkg/ssh"
	"github.com/niusmallnan/k3os/pkg/sysctl"
	"github.com/niusmallnan/k3os/pkg/util"

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
		logrus.Fatalf("failed to setting hostname: %v", err)
	}
	// setup /etc/hosts
	if err := pkgHostname.SyncHostname(); err != nil {
		logrus.Fatalf("failed to sync hostname: %v", err)
	}
	// setup ssh authorized_keys
	for _, username := range config.SSHUsers {
		if err := ssh.SetAuthorizedKeys(username, cfg); err != nil {
			logrus.Error(err)
		}
	}
	// setup kernel modules
	if err := module.LoadModules(cfg); err != nil {
		logrus.Fatalf("failed to load modules: %v", err)
	}
	// setup sysctl
	if err := sysctl.ConfigureSysctl(cfg); err != nil {
		logrus.Fatalf("failed to setting sysctl: %v", err)
	}
	// run command
	if err := command.ExecuteCommand(cfg.Runcmd); err != nil {
		logrus.Fatalf("failed to execute command: %v", err)
	}
	// run rc.local
	if err := util.RunScript("/etc/rc.local"); err != nil {
		logrus.Fatalf("failed to run rc.local: %v", err)
	}
	return nil
}
