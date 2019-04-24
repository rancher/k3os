package cliinstall

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/ghodss/yaml"

	"github.com/rancher/k3os/pkg/ask"
	"github.com/rancher/k3os/pkg/config"
)

func Run() error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	isInstall, err := ask.Ask(&cfg)
	if err != nil {
		return err
	}

	if isInstall {
		return runInstall(cfg)
	}

	cfg.K3OS.Mode = ""
	cfg.K3OS.Install = config.Install{}
	f, err := os.Create(config.SystemConfig)
	if err != nil {
		f, err = os.Create(config.LocalConfig)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(&cfg); err != nil {
		return err
	}

	f.Close()
	if err := runCCApply(); err != nil {
		return err
	}

	return exec.Command("service", "k3s", "restart").Run()
}

func runCCApply() error {
	cmd := exec.Command("/usr/sbin/ccapply", "--config")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func runInstall(cfg config.CloudConfig) error {
	var (
		err      error
		tempFile *os.File
	)

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
		cfg.K3OS.Mode = ""
		cfg.K3OS.Install = config.Install{}
		bytes, err := yaml.Marshal(&cfg)
		if err != nil {
			return err
		}
		if _, err := tempFile.Write(bytes); err != nil {
			return err
		}
		if err := tempFile.Close(); err != nil {
			return err
		}
		defer os.Remove(tempFile.Name())
	}

	cmd := exec.Command("/usr/libexec/k3os/install")
	cmd.Env = append(os.Environ(), ev...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
