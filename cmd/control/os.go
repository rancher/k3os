package control

import (
	"fmt"
	"syscall"

	"github.com/rancher/k3os/config"
	"github.com/rancher/k3os/pkg/util"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

const (
	OSUpgradeScript = "/usr/lib/k3os/k3os-upgrade"
)

type Images struct {
	Current   string  `yaml:"current,omitempty"`
	Available []Image `yaml:"available,omitempty"`
}

type Image struct {
	Version string `yaml:"version,omitempty"`
	Kernel  string `yaml:"kernel,omitempty"`
	Initrd  string `yaml:"initrd,omitempty"`
	Vmlinuz string `yaml:"vmlinuz,omitempty"`
}

func osSubcommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "upgrade",
			Usage:  "upgrade to latest version",
			Action: osUpgrade,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "image, i",
					Usage: "upgrade to a certain version",
				},
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "do not prompt for input",
				},
				cli.BoolFlag{
					Name:  "no-reboot",
					Usage: "do not reboot after upgrade",
				},
				cli.BoolFlag{
					Name:  "debug",
					Usage: "debug output",
				},
			},
		},
		{
			Name:   "list",
			Usage:  "list the current available versions",
			Action: osList,
		},
	}
}

func osUpgrade(c *cli.Context) error {
	imageVersion := c.String("image")
	rebootFlag := !c.Bool("no-reboot")
	forceFlag := c.Bool("force")
	//TODO: debug for output

	if !forceFlag && !util.PromptYes("continue with os upgrade") {
		return nil
	}
	cfg := config.LoadConfig("", false)
	upgradeURL := cfg.K3OS.Upgrade.URL
	images, err := getImages(upgradeURL)
	if err != nil {
		logrus.Fatalf("failed to get os list: %v", err)
	}

	if imageVersion == "" {
		imageVersion = images.Current
	}

	var imageObject Image
	for _, img := range images.Available {
		if imageVersion == img.Version {
			imageObject = img
			break
		}
	}

	if imageObject.Version == "" {
		logrus.Fatalf("invalid image version for %s", imageVersion)
	}

	err = util.RunScript(OSUpgradeScript, imageObject.Initrd, imageObject.Vmlinuz, imageObject.Version, imageObject.Kernel)
	if err != nil {
		logrus.Fatalf("failed to run upgrade script: %v", err)
	}

	if (rebootFlag && util.PromptYes("continue with reboot")) || forceFlag {
		syscall.Sync()
		syscall.Reboot(int(syscall.LINUX_REBOOT_CMD_RESTART))
	}

	return nil
}

func osList(c *cli.Context) error {
	cfg := config.LoadConfig("", false)
	upgradeURL := cfg.K3OS.Upgrade.URL

	images, err := getImages(upgradeURL)
	if err != nil {
		logrus.Fatalf("failed to get os list: %v", err)
	}

	fmt.Println("running:")
	fmt.Println("-", config.OSVersion)

	fmt.Println("available:")
	for i := len(images.Available) - 1; i >= 0; i-- {
		fmt.Println("-", images.Available[i].Version)
	}

	return nil
}

func getImages(upgradeURL string) (*Images, error) {
	body, err := util.HTTPLoadBytes(upgradeURL)
	if err != nil {
		return nil, err
	}

	images := &Images{}
	err = yaml.Unmarshal(body, images)
	if err != nil {
		return nil, err
	}

	return images, nil
}
