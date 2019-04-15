package network

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/niusmallnan/k3os/config"
)

func SettingProxy(cfg *config.CloudConfig) error {
	httpProxy := cfg.K3OS.Network.Proxy.HTTPProxy
	httpsProxy := cfg.K3OS.Network.Proxy.HTTPSProxy
	noProxy := cfg.K3OS.Network.Proxy.NoProxy
	proxyLines := make([]string, 0)
	if len(httpProxy) > 0 {
		for _, k := range []string{"HTTP_PROXY", "http_proxy"} {
			proxyLines = append(proxyLines, fmt.Sprintf("export %s=%s", k, httpProxy))
			if err := os.Setenv(k, httpProxy); err != nil {
				return err
			}
		}
	}
	if len(httpsProxy) > 0 {
		for _, k := range []string{"HTTPS_PROXY", "https_proxy"} {
			proxyLines = append(proxyLines, fmt.Sprintf("export %s=%s", k, httpsProxy))
			if err := os.Setenv(k, httpsProxy); err != nil {
				return err
			}
		}
	}
	if len(noProxy) > 0 {
		for _, k := range []string{"NO_PROXY", "no_proxy"} {
			proxyLines = append(proxyLines, fmt.Sprintf("export %s=%s", k, noProxy))
			if err := os.Setenv(k, noProxy); err != nil {
				return err
			}
		}
	}
	if len(proxyLines) > 0 {
		proxy := strings.Join(proxyLines, "\n")
		proxy = fmt.Sprintf("#!/bin/sh\n%s\n", proxy)
		if err := ioutil.WriteFile("/etc/profile.d/proxy.sh", []byte(proxy), 0644); err != nil {
			return err
		}
	}
	return nil
}
