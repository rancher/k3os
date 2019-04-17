package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rancher/k3os/config"
)

func ExecuteCommand(commands []config.Command) error {
	for _, cmd := range commands {
		var c *exec.Cmd
		if cmd.String != "" {
			c = exec.Command("sh", "-c", cmd.String)
		} else if len(cmd.Strings) > 0 {
			c = exec.Command(cmd.Strings[0], cmd.Strings[1:]...)
		} else {
			continue
		}
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to run %s: %v", cmd, err)
		}
	}
	return nil
}

func SetPassword(password string) error {
	if password == "" {
		return nil
	}
	cmd := exec.Command("chpasswd")
	cmd.Stdin = strings.NewReader(fmt.Sprint("rancher:", password))
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("bash", "-c", `sed -E -i 's/(rancher:.*:).*(:.*:.*:.*:.*:.*:.*)$/\1\2/' /etc/shadow`)
	return cmd.Run()
}
