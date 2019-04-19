package control

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/rancher/k3os/config"
	"github.com/rancher/k3os/pkg/util"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	GPTMBRInstallType = "gptmbr"
	MBRInstallType    = "mbr"
	EFIInstallType    = "efi"

	InstallConfigScript = "/usr/lib/k3os/k3os-install-config"
	InstallBootScript   = "/usr/lib/k3os/k3os-install-%s"
	UserConfigTempFile  = "/tmp/user_config.yml"
	EmptyConfigTempFile = "/tmp/empty_config.yml"
)

var installCommand = cli.Command{
	Name:     "install",
	Usage:    "install k3os to disk",
	HideHelp: true,
	Action:   installAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "install-type, t",
			Value: GPTMBRInstallType,
			Usage: "gptmbr, mbr, efi",
		},
		cli.StringFlag{
			Name:  "cloud-config, c",
			Usage: "cloud-config yml file - needed for SSH authorized keys",
		},
		cli.StringFlag{
			Name:  "device, d",
			Usage: "storage device",
		},
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "[ DANGEROUS! data loss can happen ] partition/format without prompting",
		},
		cli.BoolFlag{
			Name:  "no-reboot",
			Usage: "do not reboot after install",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "run installer with debug output",
		},
	},
}

func installAction(c *cli.Context) error {
	installType := c.String("install-type")
	cloudConfig := c.String("cloud-config")
	installDevice := c.String("device")
	rebootFlag := !c.Bool("no-reboot")
	forceFlag := c.Bool("force")
	//TODO: debug for output

	if installDevice == "" {
		logrus.Fatal("can not proceed without -d <dev> specified")
	}

	if cloudConfig == "" {
		logrus.Warn("cloud-config not provided")
		p, ok := util.PromptPassword()
		if !ok {
			logrus.Fatal("password and confirmation password do not match")
		}
		// create an empty file
		// TODO: direct the user to create a config file
		emptyFile, err := os.Create(EmptyConfigTempFile)
		if err != nil {
			logrus.Fatalf("failed to create empty config file, %v", err)
		}
		// set password to cloud-config
		data := make(map[interface{}]interface{}, 0)
		_, modified := util.SetValue("k3os.password", data, string(p))
		cfg := &config.CloudConfig{}
		if err := util.Convert(modified, cfg); err != nil {
			return err
		}
		if err := util.WriteToFile(modified, EmptyConfigTempFile); err != nil {
			logrus.Fatal("failed set password to cloud-config")
		}

		cloudConfig = EmptyConfigTempFile
		emptyFile.Close()
	}

	if !forceFlag && !util.PromptYes("Continue with install") {
		return nil
	}

	installBootLoader := fmt.Sprintf(InstallBootScript, installType)
	if err := util.RunScript(installBootLoader, installDevice); err != nil {
		logrus.Fatalf("failed to install boot things to disk, %v", err)
	}

	if strings.HasPrefix(cloudConfig, "http://") || strings.HasPrefix(cloudConfig, "https://") {
		if err := util.HTTPDownloadToFile(cloudConfig, UserConfigTempFile); err != nil {
			logrus.Fatalf("failed to get cloud-config via http(s): %s", cloudConfig)
		}
	} else {
		if err := util.FileCopy(cloudConfig, UserConfigTempFile); err != nil {
			logrus.Fatalf("failed to copy cloud-config: %s", cloudConfig)
		}
	}
	if err := util.RunScript(InstallConfigScript, UserConfigTempFile); err != nil {
		logrus.Fatalf("failed to install config to disk, %v", err)
	}

	if (rebootFlag && util.PromptYes("continue with reboot")) || forceFlag {
		syscall.Sync()
		syscall.Reboot(int(syscall.LINUX_REBOOT_CMD_RESTART))
	}

	return nil
}
