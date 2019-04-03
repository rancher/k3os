package network

import (
	"os"
	"os/exec"
)

func StopDhcpcd() error {
	cmd := exec.Command("dhcpcd", "-x")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
