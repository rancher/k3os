package cliinstall

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/rancher/k3os/pkg/config"
	"github.com/rancher/k3os/pkg/mode"
	"github.com/rancher/k3os/pkg/questions"
	"github.com/rancher/k3os/pkg/util"
)

func Ask(cfg *config.CloudConfig) (bool, error) {
	if ok, err := isInstall(cfg); err != nil {
		return false, err
	} else if ok {
		return true, AskInstall(cfg)
	}

	return false, AskServerAgent(cfg)
}

func isInstall(cfg *config.CloudConfig) (bool, error) {
	mode, err := mode.Get()
	if err != nil {
		return false, err
	}

	if mode == "install" {
		return true, nil
	} else if mode == "live-server" {
		return false, nil
	} else if mode == "live-agent" {
		return false, nil
	}

	i, err := questions.PromptFormattedOptions("Choose operation", 0,
		"Install to disk",
		"Configure server or agent")
	if err != nil {
		return false, err
	}

	return i == 0, nil
}

func AskInstall(cfg *config.CloudConfig) error {
	if cfg.K3OS.Install.Silent {
		return nil
	}

	if err := AskInstallDevice(cfg); err != nil {
		return err
	}

	if err := AskConfigURL(cfg); err != nil {
		return err
	}

	if cfg.K3OS.Install.ConfigURL == "" {
		if err := AskGithub(cfg); err != nil {
			return err
		}

		if err := AskPassword(cfg); err != nil {
			return err
		}

		if err := AskWifi(cfg); err != nil {
			return err
		}

		if err := AskServerAgent(cfg); err != nil {
			return err
		}
	}

	return nil
}

func AskInstallDevice(cfg *config.CloudConfig) error {
	if cfg.K3OS.Install.Device != "" {
		return nil
	}

	output, err := exec.Command("/bin/sh", "-c", "lsblk -r -o NAME,TYPE | grep -w disk | awk '{print $1}'").CombinedOutput()
	if err != nil {
		return err
	}
	fields := strings.Fields(string(output))
	i, err := questions.PromptFormattedOptions("Installation target. Device will be formatted", -1, fields...)
	if err != nil {
		return err
	}

	cfg.K3OS.Install.Device = "/dev/" + fields[i]
	return nil
}

func AskToken(cfg *config.CloudConfig, server bool) error {
	var (
		token string
		err   error
	)

	if cfg.K3OS.Token != "" {
		return nil
	}

	msg := "Token or cluster secret"
	if server {
		msg += " (optional)"
	}
	if server {
		token, err = questions.PromptOptional(msg+": ", "")
	} else {
		token, err = questions.Prompt(msg+": ", "")
	}
	cfg.K3OS.Token = token

	return err
}

func isServer(cfg *config.CloudConfig) (bool, error) {
	mode, err := mode.Get()
	if err != nil {
		return false, err
	}
	if mode == "live-server" {
		return true, nil
	} else if mode == "live-agent" || (cfg.K3OS.ServerURL != "" && cfg.K3OS.Token != "") {
		return false, nil
	}

	opts := []string{"server", "agent"}
	i, err := questions.PromptFormattedOptions("Run as server or agent?", 0, opts...)
	if err != nil {
		return false, err
	}

	return i == 0, nil
}

func AskServerAgent(cfg *config.CloudConfig) error {
	if cfg.K3OS.ServerURL != "" {
		return nil
	}

	server, err := isServer(cfg)
	if err != nil {
		return err
	}

	if server {
		return AskToken(cfg, true)
	}

	url, err := questions.Prompt("URL of server: ", "")
	if err != nil {
		return err
	}
	cfg.K3OS.ServerURL = url

	return AskToken(cfg, false)
}

func AskPassword(cfg *config.CloudConfig) error {
	if len(cfg.SSHAuthorizedKeys) > 0 || cfg.K3OS.Password != "" {
		return nil
	}

	var (
		ok   = false
		err  error
		pass string
	)

	for !ok {
		pass, ok, err = util.PromptPassword()
		if err != nil {
			return err
		}
	}

	if os.Getuid() != 0 {
		cfg.K3OS.Password = pass
		return nil
	}

	oldShadow, err := ioutil.ReadFile("/etc/shadow")
	if err != nil {
		return err
	}
	defer func() {
		ioutil.WriteFile("/etc/shadow", oldShadow, 0640)
	}()

	cmd := exec.Command("chpasswd")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("rancher:%s", pass))
	errBuffer := &bytes.Buffer{}
	cmd.Stdout = os.Stdout
	cmd.Stderr = errBuffer

	if err := cmd.Run(); err != nil {
		os.Stderr.Write(errBuffer.Bytes())
		return err
	}

	f, err := os.Open("/etc/shadow")
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ":")
		if len(fields) > 1 && fields[0] == "rancher" {
			cfg.K3OS.Password = fields[1]
			return nil
		}
	}

	return scanner.Err()
}

func AskWifi(cfg *config.CloudConfig) error {
	if len(cfg.K3OS.Wifi) > 0 {
		return nil
	}

	ok, err := questions.PromptBool("Configure WiFi?", false)
	if !ok || err != nil {
		return err
	}

	for {
		name, err := questions.Prompt("WiFi Name: ", "")
		if err != nil {
			return err
		}

		pass, err := questions.Prompt("WiFi Passphrase: ", "")
		if err != nil {
			return err
		}

		cfg.K3OS.Wifi = append(cfg.K3OS.Wifi, config.Wifi{
			Name:       name,
			Passphrase: pass,
		})

		ok, err := questions.PromptBool("Configure another WiFi network?", false)
		if !ok || err != nil {
			return err
		}
	}
}

func AskGithub(cfg *config.CloudConfig) error {
	if len(cfg.SSHAuthorizedKeys) > 0 || cfg.K3OS.Password != "" {
		return nil
	}

	ok, err := questions.PromptBool("Authorize GitHub users to SSH?", false)
	if !ok || err != nil {
		return err
	}

	str, err := questions.Prompt("Comma separated list of GitHub users to authorize: ", "")
	if err != nil {
		return err
	}

	for _, s := range strings.Split(str, ",") {
		cfg.SSHAuthorizedKeys = append(cfg.SSHAuthorizedKeys, "github:"+strings.TrimSpace(s))
	}

	return nil
}

func AskConfigURL(cfg *config.CloudConfig) error {
	if cfg.K3OS.Install.ConfigURL != "" {
		return nil
	}

	ok, err := questions.PromptBool("Config system with cloud-init file?", false)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	str, err := questions.Prompt("cloud-init file location (file path or http URL): ", "")
	if err != nil {
		return err
	}

	cfg.K3OS.Install.ConfigURL = str
	return nil
}
