package cc

import (
	"github.com/rancher/k3os/pkg/config"
	"github.com/urfave/cli"
)

type applier func(cfg *config.CloudConfig) error

func runApplies(cfg *config.CloudConfig, appliers ...applier) error {
	var errors []error

	for _, a := range appliers {
		err := a(cfg)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return cli.NewMultiError(errors...)
	}

	return nil
}

func RunApply(cfg *config.CloudConfig) error {
	return runApplies(cfg,
		ApplySSHKeysWithNet,
		ApplyWriteFiles,
		ApplyEnvironment,
		ApplyRuncmd,
		ApplyInstall,
		ApplyK3SInstall,
	)
}

func ConfigApply(cfg *config.CloudConfig) error {
	return runApplies(cfg,
		ApplyK3SWithRestart,
	)
}

func BootApply(cfg *config.CloudConfig) error {
	return runApplies(cfg,
		ApplyDataSource,
		ApplyModules,
		ApplySysctls,
		ApplyHostname,
		ApplyDNS,
		ApplyWifi,
		ApplyPassword,
		ApplySSHKeys,
		ApplyK3SNoRestart,
		ApplyWriteFiles,
		ApplyEnvironment,
		ApplyBootcmd,
	)
}

func InitApply(cfg *config.CloudConfig) error {
	return runApplies(cfg,
		ApplyModules,
		ApplySysctls,
		ApplyHostname,
		ApplyWriteFiles,
		ApplyEnvironment,
		ApplyInitcmd,
	)
}
