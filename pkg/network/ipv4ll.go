package network

import (
	"encoding/binary"
	"errors"
	"math/rand"
	"net"

	"github.com/j-keck/arping"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func SettingLocalLinkIP(link netlink.Link) error {
	name := link.Attrs().Name
	addresses, err := getAddresses(link)
	if err != nil {
		return err
	}
	for _, addr := range addresses {
		if addr.String()[:7] == "169.254" {
			logrus.Infof("local-link ip already set on interface %s", name)
			return nil
		}
	}
	g, err := pseudoRandomGenerator(link.Attrs().HardwareAddr)
	if err != nil {
		return err
	}
	for i := 0; i < 10; i++ {
		rg := rand.New(*g)
		rn := rg.Uint32()
		ip := generateIPV4LL(rn)
		if ip[2] == 0 || ip[2] == 255 {
			i--
			continue
		}
		_, _, err := arping.PingOverIfaceByName(ip, name)
		if err != nil {
			// this ip is not being used
			addr, err := netlink.ParseAddr(ip.String() + "/16")
			if err != nil {
				return err
			}
			if err := netlink.AddrAdd(link, addr); err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("couldn't find a suitable ipv4ll")
}

func RemoveLocalLinkIP(link netlink.Link) error {
	addresses, err := getAddresses(link)
	if err != nil {
		return err
	}
	for _, addr := range addresses {
		if addr.String()[:7] == "169.254" {
			if err := netlink.AddrDel(link, &addr); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func generateIPV4LL(random uint32) net.IP {
	byte1 := random & 255 // use least significant 8 bits
	byte2 := random >> 24 // use most significant 8 bits
	return []byte{169, 254, byte(byte1), byte(byte2)}
}

func pseudoRandomGenerator(haAddr []byte) (*rand.Source, error) {
	seed, _ := binary.Varint(haAddr)
	src := rand.NewSource(seed)
	return &src, nil
}
