package apply

import (
	"github.com/rancher/k3os/pkg/command"
	"github.com/rancher/k3os/pkg/config"
	"github.com/rancher/k3os/pkg/hostname"
	"github.com/rancher/k3os/pkg/module"
	"github.com/rancher/k3os/pkg/sysctl"
	"github.com/rancher/k3os/pkg/writefile"
)

func ApplyModules(cfg *config.CloudConfig) error {
	return module.LoadModules(cfg)
}

func ApplySysctls(cfg *config.CloudConfig) error {
	return sysctl.ConfigureSysctl(cfg)
}

func ApplyHostname(cfg *config.CloudConfig) error {
	return hostname.SetHostname(cfg)
}

func ApplyPassword(cfg *config.CloudConfig) error {
	return command.SetPassword(cfg.K3OS.Password)
}

func ApplyRuncmd(cfg *config.CloudConfig) error {
	return command.ExecuteCommand(cfg.Runcmd)
}

func ApplyBootcmd(cfg *config.CloudConfig) error {
	return command.ExecuteCommand(cfg.Bootcmd)
}

func ApplyInitcmd(cfg *config.CloudConfig) error {
	return command.ExecuteCommand(cfg.Initcmd)
}

func ApplyWriteFiles(cfg *config.CloudConfig) error {
	writefile.WriteFiles(cfg)
	return nil
}

func ApplyK3S(cfg *config.CloudConfig) error {
	buf :=
	contents =

}
