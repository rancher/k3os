package control

import (
	"fmt"
	"os"

	"github.com/niusmallnan/k3os/config"

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
	app.Usage = fmt.Sprintf("Control and configure K3OS(%s)", config.OSBuildDate)
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
			Name:            "entrypoint",
			Hidden:          true,
			HideHelp:        true,
			SkipFlagParsing: true,
			Action:          entryPoint,
		},
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
	cfg := config.LoadConfig("", false)
	n := fmt.Sprintf("%s:%s", cfg.K3OS.Upgrade.Image, config.OSVersion)
	fmt.Fprintf(c.App.Writer, "version %s from k3os image %s\n", c.App.Version, n)
}
