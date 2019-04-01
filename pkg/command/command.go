package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/niusmallnan/k3os/config"
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
