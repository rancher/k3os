package apply

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

func ApplyK3S(cfg *config.CloudConfig) error {
	buf := &bytes.Buffer{}
	if cfg.K3OS.ServerURL != "" {
		buf.WriteString(fmt.Sprintf("export K3S_URL=\"%s\"\n", cfg.K3OS.ServerURL))
	}

	if strings.HasPrefix(cfg.K3OS.Token, "K10") {
		buf.WriteString(fmt.Sprintf("export K3S_TOKEN=\"%s\"\n", cfg.K3OS.Token))
	} else if cfg.K3OS.Token != "" {
		buf.WriteString(fmt.Sprintf("export K3S_CLUSTER_SECRET=\"%s\"\n", cfg.K3OS.Token))
	}

	if buf.Len() > 0 {
		if err := os.MkdirAll("/etc/rancher/k3s", 755); err != nil {
			return err
		}
		return ioutil.WriteFile("/etc/rancher/k3s/k3s.env", buf.Bytes(), 644)
	}

	return nil
}
