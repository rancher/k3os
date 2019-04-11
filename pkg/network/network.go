package network

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"syscall"

	"github.com/niusmallnan/k3os/config"

	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
)

func SettingNetwork(cfg *config.CloudConfig) error {
	interfaces := cfg.K3OS.Network.Interfaces
	links, err := getLinkList()
	if err != nil {
		return err
	}
	exist, err := CheckDhcpcd()
	if err != nil {
		return err
	}
	// dhcp mode on all interfaces
	if len(interfaces) <= 0 {
		if !exist {
			// cleanup network settings
			if err := cleanupSettings(links); err != nil {
				return err
			}
			// start dhcpcd process
			if err := StartDhcpcd([]string{}); err != nil {
				return fmt.Errorf("failed to start dhcpcd: %v", err)
			}
		}
		logrus.Infoln("all link will use dhcp mode")
	} else {
		// find links which is dhcp mode
		dLinks := findDhcpLinks(interfaces, links)
		// no interface use dhcp
		if len(dLinks) <= 0 {
			if err := ReleaseDhcpcd(""); err != nil {
				// consider no running dhcpcd process situation, so ignore the error
				logrus.Warnf("failed to release dhcpcd, network setting will be continue: %v", err)
			}
			reCheck, err := CheckDhcpcd()
			if err != nil {
				return err
			}
			if reCheck {
				return errors.New("dhcpcd process is still running")
			}
		} else {
			reCheck, err := CheckDhcpcd()
			if err != nil {
				return err
			}
			if !reCheck {
				// make sure dhcpcd process exist when execute k3os-netinit with no reboot
				// -w wait for an address to be assigned before forking to the background
				if err := StartDhcpcd([]string{"-w"}); err != nil {
					return err
				}
			}
			// release dhcp for interface which not use dhcp mode
			for k := range interfaces {
				link, err := netlink.LinkByName(k)
				if err != nil || link == nil {
					logrus.Warnf("interface %s not exist", k)
					continue
				}
				if err := ReleaseDhcpcd(k); err != nil {
					// consider no running dhcpcd process situation, so ignore the error
					logrus.Warnf("failed to release dhcp on link %s , network setting will be continue: %v", k, err)
				}
			}
		}
		// cleanup network settings
		if err := cleanupSettings(links); err != nil {
			return err
		}
		// setting cloud-config address for interface
		if err := applySettings(interfaces); err != nil {
			return err
		}
		// setting dhcp address for interface
		for _, link := range dLinks {
			logrus.Infof("link %s will use dhcp mode", link.Attrs().Name)
			if err := RequestDhcpcd(link.Attrs().Name); err != nil {
				return err
			}
		}
	}
	return nil
}

func applySettings(interfaces map[string]config.InterfaceConfig) error {
	for k, v := range interfaces {
		link, err := netlink.LinkByName(k)
		if err != nil || link == nil {
			logrus.Warnf("interface %s not exist", k)
			continue
		}
		if len(v.Addresses) <= 0 {
			return fmt.Errorf("interface %s addresses property length is 0", k)
		}
		for _, addr := range v.Addresses {
			// setting address to specific interface
			if err := settingAddress(link, addr); err != nil {
				return err
			}
		}
		// setting default gateway with metric property
		if err := settingGateway(v.Gateway, v.Metric, link); err != nil {
			return err
		}
	}
	return nil
}

func cleanupSettings(links []netlink.Link) error {
	// remove useless default gateway
	if err := removeGateway(); err != nil {
		return err
	}
	for _, link := range links {
		// remove all useless address on specific interface
		if err := removeAddresses(link); err != nil {
			return err
		}
	}
	return nil
}

func findDhcpLinks(interfaces map[string]config.InterfaceConfig, links []netlink.Link) []netlink.Link {
	dhcpLinks := make([]netlink.Link, 0)
	for _, link := range links {
		if _, ok := interfaces[link.Attrs().Name]; !ok {
			dhcpLinks = append(dhcpLinks, link)
		}
	}
	return dhcpLinks
}

func getLinkList() ([]netlink.Link, error) {
	var valid []netlink.Link
	links, err := netlink.LinkList()
	if err != nil {
		return valid, err
	}
	for _, l := range links {
		name := l.Attrs().Name
		if name == "lo" || strings.Contains(name, "flannel") || strings.Contains(name, "veth") {
			continue
		}
		valid = append(valid, l)
	}
	return valid, nil
}

func removeAddresses(link netlink.Link) error {
	addresses, err := netlink.AddrList(link, nl.FAMILY_ALL)
	if err != nil {
		return err
	}
	for _, addr := range addresses {
		if err := netlink.AddrDel(link, &addr); err != nil && err != syscall.EEXIST {
			return err
		}
	}
	return nil
}

func removeGateway() error {
	routes, err := netlink.RouteList(nil, nl.FAMILY_ALL)
	if err != nil {
		return err
	}
	for _, r := range routes {
		if r.Dst == nil {
			if err := netlink.RouteDel(&r); err != nil {
				return err
			}
		}
	}
	return nil
}

func settingAddress(link netlink.Link, address string) error {
	addr, err := netlink.ParseAddr(address)
	if err != nil {
		return err
	}
	if err := netlink.AddrAdd(link, addr); err != nil && err != syscall.EEXIST {
		return err
	}
	return nil
}

func settingGateway(gateway string, metric int, link netlink.Link) error {
	if gateway == "" {
		logrus.Warnf("interface %s's gateway property is not setting", link.Attrs().Name)
		return nil
	}
	gw := net.ParseIP(gateway)
	if gw == nil {
		return errors.New("invalid gateway address: " + gateway)
	}
	// Metrics are used to prefer an interface over another one, lowest wins.
	//  Dhcpcd will supply a default metric	of 200 + if_nametoindex(3).
	//  An extra 100 will be added for wireless interfaces.
	// Reference: dhcpcd.conf.5
	priority := 400 + link.Attrs().Index
	if metric > 0 {
		priority = metric
	}
	route := netlink.Route{
		Scope: netlink.SCOPE_UNIVERSE,
		Gw:    gw,
		// cloud-config user-defined gateway begin with 400 + if_nametoindex(3)
		Priority: priority,
	}
	if err := netlink.RouteAdd(&route); err != nil && err != syscall.EEXIST {
		return err
	}
	return nil
}
