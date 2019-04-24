package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ExecuteCommand(commands []string) error {
	for _, cmd := range commands {
		c := exec.Command("sh", "-c", cmd)
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
	if strings.HasPrefix(password, "$") {
		cmd.Args = append(cmd.Args, "-e")
	}
	cmd.Stdin = strings.NewReader(fmt.Sprint("rancher:", password))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
