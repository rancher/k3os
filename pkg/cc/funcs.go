package cc

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

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
	return ApplyK3S(cfg, true)
}

func ApplyK3SNoRestart(cfg *config.CloudConfig) error {
	return ApplyK3S(cfg, false)
}

func ApplyK3S(cfg *config.CloudConfig, restart bool) error {
	var args []string
	vars := []string{
		"INSTALL_K3S_NAME=service",
		"INSTALL_K3S_SKIP_DOWNLOAD=true",
		"INSTALL_K3S_BIN_DIR=/sbin",
		"INSTALL_K3S_BIN_DIR_READ_ONLY=true",
	}

	if !restart {
		vars = append(vars, "INSTALL_K3S_SKIP_START=true")
	}

	if cfg.K3OS.ServerURL != "" {
		vars = append(vars, fmt.Sprintf("K3S_URL=\"%s\"\n", cfg.K3OS.ServerURL))
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
