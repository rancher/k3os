package cc

import (
	"reflect"
	"runtime"

	"github.com/rancher/k3os/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type applier func(cfg *config.CloudConfig) error

func runApplies(cfg *config.CloudConfig, appliers ...applier) error {
	var errors []error

	if l := logrus.GetLevel(); l >= logrus.DebugLevel {
		c := make([]uintptr, 2)
		n := runtime.Callers(2, c)
		s := runtime.CallersFrames(c[:n])
		f, _ := s.Next()
		logrus.Debugf(">>> %s", f.Function)
		defer logrus.Debugf("<<< %s", f.Function)
	}

	for _, a := range appliers {
		logrus.Debugf("+++ %s", runtime.FuncForPC(reflect.ValueOf(a).Pointer()).Name())
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
