package network

import (
	"errors"
	"fmt"
	"net"
	"syscall"

	"github.com/niusmallnan/k3os/config"

	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
)

var defaultMetric = 400

func SettingNetwork(cfg *config.CloudConfig) error {
	interfaces := cfg.K3OS.Network.Interfaces
	if len(interfaces) > 0 {
		if err := ReleaseDhcpcd(""); err != nil {
			// consider no running dhcpcd process situation, so ignore the error
			logrus.Errorf("failed to release dhcpcd, network setting will be continue: %v", err)
		}
		if err := CheckDhcpcd(); err != nil {
			return err
		}
		// remove useless default gateway
		if err := removeGateway(); err != nil {
			return err
		}
		for k, v := range interfaces {
			link, err := netlink.LinkByName(k)
			if err != nil || link == nil {
				logrus.Errorf("interface %s not exist", k)
				continue
			}
			// remove all useless address on specific interface
			if err := removeAddresses(link); err != nil {
				return err
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
			if err := settingGateway(v.Gateway, link); err != nil {
				return err
			}
		}
	} else {
		// start dhcpcd process, ignore the prev-settings which will be clean after reboot
		if err := StartDhcpcd(); err != nil {
			return fmt.Errorf("failed to start dhcpcd: %v", err)
		}
	}
	return nil
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

func settingGateway(gateway string, link netlink.Link) error {
	if gateway == "" {
		return errors.New("gateway can not be empty")
	}
	gw := net.ParseIP(gateway)
	if gw == nil {
		return errors.New("invalid gateway address: " + gateway)
	}

	// Metrics are used to prefer an interface over another one, lowest wins.
	//  Dhcpcd will supply a default metric	of 200 + if_nametoindex(3).
	//  An extra 100 will be added for wireless interfaces.
	// Reference: dhcpcd.conf.5
	route := netlink.Route{
		Scope:    netlink.SCOPE_UNIVERSE,
		Gw:       gw,
		Priority: defaultMetric + link.Attrs().Index,
	}

	if err := netlink.RouteAdd(&route); err != nil && err != syscall.EEXIST {
		return err
	}
	return nil
}
