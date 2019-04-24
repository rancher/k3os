package ask

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/rancher/k3os/pkg/config"
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
	if cfg.K3OS.Mode == "install" {
		return true, nil
	} else if cfg.K3OS.Mode == "live-server" {
		return false, nil
	} else if cfg.K3OS.Mode == "live-agent" {
		return false, nil
	}

	i, err := questions.PromptFormattedOptions("Choose operation", -1,
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

	if err := AskInstallEFI(cfg); err != nil {
		return err
	}

	if !cfg.K3OS.Install.EFI {
		if err := AskMSDOS(cfg); err != nil {
			return err
		}
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

		if err := AskServerAgent(cfg); err != nil {
			return err
		}
	}

	return nil
}

func AskMSDOS(cfg *config.CloudConfig) error {
	if cfg.K3OS.Install.MSDOS {
		return nil
	}

	i, err := questions.PromptFormattedOptions("Choose installation partition table type", 0,
		"gpt",
		"msdos")
	if err != nil {
		return err
	}

	cfg.K3OS.Install.MSDOS = i == 1
	return nil
}

func AskInstallDevice(cfg *config.CloudConfig) error {
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

func AskInstallEFI(cfg *config.CloudConfig) error {
	if cfg.K3OS.Install.EFI {
		return nil
	}

	if _, err := os.Stat("/sys/firmware/efi"); err != nil {
		return nil
	}

	cfg.K3OS.Install.EFI = true
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
	if cfg.K3OS.Mode == "live-server" {
		return true, nil
	} else if cfg.K3OS.Mode == "live-agent" {
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
	if len(cfg.SSHAuthorizedKeys) > 0 {
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
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

func AskGithub(cfg *config.CloudConfig) error {
	ok, err := questions.PromptBool("Authorize GitHub users to SSH?", false)
	if !ok || err != nil {
		return err
	}

	str, err := questions.Prompt("Comma seperated list of GitHub users or organizations to authorize: ", "")
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
