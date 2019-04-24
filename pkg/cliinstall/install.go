package cliinstall

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/rancher/k3os/pkg/ask"
	"github.com/rancher/k3os/pkg/config"
)

func Run(args []string) error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	if err := ask.AskInstall(&cfg); err != nil {
		return err
	}

	var tempFile *os.File

	if cfg.K3OS.Install.ConfigURL == "" {
		tempFile, err = ioutil.TempFile("/tmp", "k3os.XXXXXXXX")
		if err != nil {
			return err
		}
		defer tempFile.Close()

		cfg.K3OS.Install.ConfigURL = tempFile.Name()
	}

	ev, err := config.ToEnv(cfg)
	if err != nil {
		return err
	}

	if tempFile != nil {
		cfg.K3OS.Install = config.Install{}
		if err := json.NewEncoder(tempFile).Encode(&cfg); err != nil {
			return err
		}
		if err := tempFile.Close(); err != nil {
			return err
		}
		defer os.Remove(tempFile.Name())
	}

	if len(args) > 0 {
		cmd := exec.Command(args[0])
		cmd.Env = append(os.Environ(), ev...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}

	return nil
}
