package network

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rancher/k3os/config"

	"github.com/sirupsen/logrus"
)

const (
	proxyProfile = "/etc/profile.d/proxy.sh"
)

func SettingProxy() error {
	cfg := config.LoadConfig("", false)
	httpProxy := cfg.K3OS.Network.Proxy.HTTPProxy
	httpsProxy := cfg.K3OS.Network.Proxy.HTTPSProxy
	noProxy := cfg.K3OS.Network.Proxy.NoProxy
	proxyLines := make([]string, 0)
	if len(httpProxy) > 0 {
		for _, k := range []string{"HTTP_PROXY", "http_proxy"} {
			if v, ok := cfg.K3OS.Environment[k]; ok {
				proxyLines = append(proxyLines, fmt.Sprintf("export %s=%s", k, v))
			}
		}
	}
	if len(httpsProxy) > 0 {
		for _, k := range []string{"HTTPS_PROXY", "https_proxy"} {
			if v, ok := cfg.K3OS.Environment[k]; ok {
				proxyLines = append(proxyLines, fmt.Sprintf("export %s=%s", k, v))
			}
		}
	}
	if len(noProxy) > 0 {
		for _, k := range []string{"NO_PROXY", "no_proxy"} {
			if v, ok := cfg.K3OS.Environment[k]; ok {
				proxyLines = append(proxyLines, fmt.Sprintf("export %s=%s", k, v))
			}
		}
	}
	if len(proxyLines) > 0 {
		proxy := strings.Join(proxyLines, "\n")
		proxy = fmt.Sprintf("#!/bin/sh\n%s\n", proxy)
		if err := ioutil.WriteFile(proxyProfile, []byte(proxy), 0644); err != nil {
			return err
		}
	}
	return nil
}

func SettingProxyEnvironments(cfg *config.CloudConfig) {
	httpProxy := cfg.K3OS.Network.Proxy.HTTPProxy
	httpsProxy := cfg.K3OS.Network.Proxy.HTTPSProxy
	noProxy := cfg.K3OS.Network.Proxy.NoProxy
	if httpProxy != "" {
		if err := os.Setenv("HTTP_PROXY", httpProxy); err != nil {
			logrus.Errorf("unable to set HTTP_PROXY environment: %v", err)
		}
		if err := config.Set("k3os.environment.http_proxy", httpProxy); err != nil {
			logrus.Errorf("unable to set k3os.environment.http_proxy: %v", err)
		}
		if err := config.Set("k3os.environment.HTTP_PROXY", httpProxy); err != nil {
			logrus.Errorf("unable to set k3os.environment.HTTP_PROXY: %v", err)
		}
	}
	if httpsProxy != "" {
		if err := os.Setenv("HTTPS_PROXY", httpsProxy); err != nil {
			logrus.Errorf("unable to set HTTPS_PROXY environment: %v", err)
		}
		if err := config.Set("k3os.environment.https_proxy", httpsProxy); err != nil {
			logrus.Errorf("unable to set k3os.environment.https_proxy: %v", err)
		}
		if err := config.Set("k3os.environment.HTTPS_PROXY", httpsProxy); err != nil {
			logrus.Errorf("unable to set k3os.environment.HTTPS_PROXY: %v", err)
		}
	}
	if noProxy != "" {
		if err := os.Setenv("NO_PROXY", noProxy); err != nil {
			logrus.Errorf("unable to set NO_PROXY environment: %v", err)
		}
		if err := config.Set("k3os.environment.no_proxy", noProxy); err != nil {
			logrus.Errorf("unable to set k3os.environment.no_proxy: %v", err)
		}
		if err := config.Set("k3os.environment.NO_PROXY", noProxy); err != nil {
			logrus.Errorf("unable to set k3os.environment.NO_PROXY: %v", err)
		}
	}
}
