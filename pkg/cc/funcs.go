package cc

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/rancher/k3os/pkg/mode"

	"github.com/sirupsen/logrus"

	"github.com/rancher/k3os/pkg/command"
	"github.com/rancher/k3os/pkg/config"
	"github.com/rancher/k3os/pkg/hostname"
	"github.com/rancher/k3os/pkg/module"
	"github.com/rancher/k3os/pkg/ssh"
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

func ApplySSHKeys(cfg *config.CloudConfig) error {
	return ssh.SetAuthorizedKeys(cfg, false)
}

func ApplySSHKeysWithNet(cfg *config.CloudConfig) error {
	return ssh.SetAuthorizedKeys(cfg, true)
}

func ApplyK3SWithRestart(cfg *config.CloudConfig) error {
	return ApplyK3S(cfg, true, false)
}

func ApplyK3SInstall(cfg *config.CloudConfig) error {
	return ApplyK3S(cfg, true, true)
}

func ApplyK3SNoRestart(cfg *config.CloudConfig) error {
	return ApplyK3S(cfg, false, false)
}

func ApplyK3S(cfg *config.CloudConfig, restart, install bool) error {
	mode, err := mode.Get()
	if err != nil {
		return err
	}
	if mode == "install" {
		return nil
	}

	k3sExists := false
	k3sLocalExists := false
	if _, err := os.Stat("/sbin/k3s"); err == nil {
		k3sExists = true
	}
	if _, err := os.Stat("/usr/local/bin/k3s"); err == nil {
		k3sLocalExists = true
	}

	args := cfg.K3OS.K3sArgs
	vars := []string{
		"INSTALL_K3S_NAME=service",
	}

	if !k3sExists && !restart {
		return nil
	}

	if k3sExists {
		vars = append(vars, "INSTALL_K3S_SKIP_DOWNLOAD=true")
		vars = append(vars, "INSTALL_K3S_BIN_DIR=/sbin")
		vars = append(vars, "INSTALL_K3S_BIN_DIR_READ_ONLY=true")
	} else if k3sLocalExists {
		vars = append(vars, "INSTALL_K3S_SKIP_DOWNLOAD=true")
	} else if !install {
		return nil
	}

	if !restart {
		vars = append(vars, "INSTALL_K3S_SKIP_START=true")
	}

	if cfg.K3OS.ServerURL == "" {
		if len(args) == 0 {
			args = append(args, "server")
		}
	} else {
		vars = append(vars, fmt.Sprintf("K3S_URL=\"%s\"\n", cfg.K3OS.ServerURL))
		if len(args) == 0 {
			args = append(args, "agent")
		}
	}

	if strings.HasPrefix(cfg.K3OS.Token, "K10") {
		vars = append(vars, fmt.Sprintf("K3S_TOKEN=\"%s\"\n", cfg.K3OS.Token))
	} else if cfg.K3OS.Token != "" {
		vars = append(vars, fmt.Sprintf("K3S_CLUSTER_SECRET=\"%s\"\n", cfg.K3OS.Token))
	}

	var labels []string
	for k, v := range cfg.K3OS.Labels {
		labels = append(labels, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(labels)

	for _, l := range labels {
		args = append(args, "--kubelet-arg", "node-labels="+l)
	}

	for _, taint := range cfg.K3OS.Taints {
		args = append(args, "--kubelet-arg", "register-with-taints="+taint)
	}

	cmd := exec.Command("/usr/libexec/k3os/k3s-install.sh", args...)
	cmd.Env = append(os.Environ(), vars...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	logrus.Debugf("Running %s %v %v", cmd.Path, cmd.Args, vars)

	return cmd.Run()
}

func ApplyInstall(cfg *config.CloudConfig) error {
	mode, err := mode.Get()
	if err != nil {
		return err
	}
	if mode != "install" {
		return nil
	}

	cmd := exec.Command("os-config")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
