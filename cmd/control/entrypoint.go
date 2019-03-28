package control

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/niusmallnan/k3os/config"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func entryPoint(c *cli.Context) error {
	setupNecessaryFs()
	if len(os.Args) < 3 {
		return nil
	}
	binary, err := exec.LookPath(os.Args[2])
	if err != nil {
		return err
	}
	return syscall.Exec(binary, os.Args[2:], os.Environ())
}

func setupNecessaryFs() {
	if _, err := os.Stat(config.CloudConfigDir); err != nil && os.IsNotExist(err) {
		err := os.Mkdir(config.CloudConfigDir, 0644)
		if err != nil {
			logrus.Error(err)
		}
	}
}
