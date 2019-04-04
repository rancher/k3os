package network

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/niusmallnan/k3os/config"
)

func SettingProxy(cfg *config.CloudConfig) error {
	address := cfg.K3OS.Network.Proxy.Address
	noProxy := cfg.K3OS.Network.Proxy.NoProxy
	protocol := cfg.K3OS.Network.Proxy.Protocol
	if address != "" && protocol != "" {
		proxyLines := make([]string, 0)
		switch strings.ToLower(protocol) {
		case "http":
			for _, k := range []string{"HTTP_PROXY", "http_proxy"} {
				proxyLines = append(proxyLines, fmt.Sprintf("export %s=%s", k, address))
				if err := os.Setenv(k, address); err != nil {
					return err
				}
			}
			break
		case "https":
			for _, k := range []string{"HTTPS_PROXY", "https_proxy"} {
				proxyLines = append(proxyLines, fmt.Sprintf("export %s=%s", k, address))
				if err := os.Setenv(k, address); err != nil {
					return err
				}
			}
			break
		default:
			return errors.New("proxy's protocol not match http or https")
		}
		if len(noProxy) > 0 {
			for _, k := range []string{"NO_PROXY", "no_proxy"} {
				proxyLines = append(proxyLines, fmt.Sprintf("export %s=%s", k, noProxy))
				if err := os.Setenv(k, address); err != nil {
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
	}
	return nil
}
