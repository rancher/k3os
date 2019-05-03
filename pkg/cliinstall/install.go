package cliinstall

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/ghodss/yaml"
	"github.com/rancher/k3os/pkg/config"
	"github.com/rancher/k3os/pkg/questions"
)

func Run() error {
	fmt.Println("\nRunning k3OS configuration")

	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	isInstall, err := Ask(&cfg)
	if err != nil {
		return err
	}

	if isInstall {
		return runInstall(cfg)
	}

	bytes, err := config.ToBytes(cfg)
	if err != nil {
		return err
	}

	f, err := os.Create(config.SystemConfig)
	if err != nil {
		f, err = os.Create(config.LocalConfig)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	if _, err := f.Write(bytes); err != nil {
		return err
	}

	f.Close()
	return runCCApply()
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

	installBytes, err := config.PrintInstall(cfg)
	if err != nil {
		return err
	}

	val, err := questions.PromptBool("\nConfiguration\n"+"-------------\n\n"+
		string(installBytes)+
		"\nYour disk will be formatted and k3OS will be installed with the above configuration.\nContinue?", false)
	if err != nil || !val {
		return err
	}

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
		cfg.K3OS.Install = nil
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
