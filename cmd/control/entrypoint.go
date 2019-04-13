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
	if _, err := os.Stat(config.CloudConfigDir); os.IsNotExist(err) {
		err := os.MkdirAll(config.CloudConfigDir, 0755)
		if err != nil {
			logrus.Error(err)
		}
	} else if err != nil {
		logrus.Error(err)
	}
}
