package environment

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rancher/k3os/config"
)

const (
	K3SProfile  = "/etc/profile.d/k3s.sh"
	K3OSProfile = "/etc/profile.d/k3os.sh"
)

func SettingEnvironments() error {
	cfg := config.LoadConfig("", false)
	k3sLines := make([]string, 0)
	k3osLines := make([]string, 0)
	for k, v := range cfg.K3OS.Environment {
		lower := strings.ToLower(k)
		// ignore network proxy settings because of net-init will setting this environments
		if "http_proxy" == lower && cfg.K3OS.Network.Proxy.HTTPProxy != "" ||
			"https_proxy" == lower && cfg.K3OS.Network.Proxy.HTTPSProxy != "" ||
			"no_proxy" == lower && cfg.K3OS.Network.Proxy.NoProxy != "" {
			continue
		}
		if strings.Contains(lower, "k3s") {
			k3sLines = append(k3sLines, fmt.Sprintf("export %s=%s", k, v))
		} else {
			k3osLines = append(k3osLines, fmt.Sprintf("export %s=%s", k, v))
		}
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	if len(k3sLines) > 0 {
		k3s := strings.Join(k3sLines, "\n")
		k3s = fmt.Sprintf("#!/bin/sh\n%s\n", k3s)
		if err := ioutil.WriteFile(K3SProfile, []byte(k3s), 0644); err != nil {
			return err
		}
	}
	if len(k3osLines) > 0 {
		k3os := strings.Join(k3osLines, "\n")
		k3os = fmt.Sprintf("#!/bin/sh\n%s\n", k3os)
		if err := ioutil.WriteFile(K3OSProfile, []byte(k3os), 0644); err != nil {
			return err
		}
	}
	return nil
}
