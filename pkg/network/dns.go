package network

import (
	"github.com/rancher/k3os/config"

	"github.com/docker/libnetwork/resolvconf"
)

func SettingDNS(cfg *config.CloudConfig) error {
	servers := cfg.K3OS.Network.DNS.Nameservers
	searches := cfg.K3OS.Network.DNS.Searches
	if len(servers) > 0 || len(searches) > 0 {
		// TODO: dhcpcd process will be replace /etc/resolv.conf, let's deal with that later
		if _, err := resolvconf.Build("/etc/resolv.conf", servers, searches, nil); err != nil {
			return err
		}
	}
	return nil
}
