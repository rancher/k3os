package hostname

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/niusmallnan/k3os/config"
)

func SetHostname(c *config.CloudConfig) error {
	hostname := c.Hostname
	if hostname == "" {
		return nil
	}
	return syscall.Sethostname([]byte(hostname))
}

func SyncHostname() error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	if hostname == "" {
		return nil
	}
	hosts, err := os.Open("/etc/hosts")
	defer hosts.Close()
	if err != nil {
		return err
	}
	lines := bufio.NewScanner(hosts)
	content := ""
	for lines.Scan() {
		line := strings.TrimSpace(lines.Text())
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == "127.0.1.1" {
			content += "127.0.1.1 " + hostname + "\n"
			continue
		}
		content += line + "\n"
	}
	return ioutil.WriteFile("/etc/hosts", []byte(content), 0600)
}
