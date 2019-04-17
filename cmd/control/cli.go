package control

import (
	"fmt"
	"os"

	"github.com/rancher/k3os/config"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func init() {
	cli.VersionPrinter = versionPrinter
}

func Main() {
	// TODO: rsyslog need to be added here.
	app := cli.NewApp()
	app.Author = "Rancher Labs, Inc."
	app.Before = beforeFunc
	app.EnableBashCompletion = true
	app.Name = os.Args[0]
	app.Usage = fmt.Sprintf("control and configure K3OS(%s)", config.OSBuildDate)
	app.Version = config.OSVersion
	app.Commands = []cli.Command{
		{
			Name:        "config",
			ShortName:   "c",
			Usage:       "configure settings",
			HideHelp:    true,
			Subcommands: configCommands(),
		},
		{
			Name:        "os",
			Usage:       "operating system upgrade/downgrade",
			HideHelp:    true,
			Subcommands: osSubcommands(),
		},
		installCommand,
	}
	app.Run(os.Args)
}

func beforeFunc(c *cli.Context) error {
	if os.Getuid() != 0 {
		logrus.Fatalf("%s: need to be root", os.Args[0])
	}
	return nil
}

func versionPrinter(c *cli.Context) {
	// TODO: print version and checksum for cli, kernel, rootfs
	fmt.Fprintf(c.App.Writer, config.OSVersion)
}
