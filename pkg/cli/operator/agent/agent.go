package agent

import (
	"context"
	"fmt"
	"os"

	"github.com/rancher/k3os/pkg/controller"
	"github.com/rancher/k3os/pkg/system"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Command is the `agent` sub-command, it is the k3OS resource controller.
func Command() cli.Command {
	return cli.Command{
		Name:  "agent",
		Usage: "custom resource controller(s)",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "name",
				EnvVar:   "K3OS_OPERATOR_NAME",
				Required: true,
			},
			cli.StringFlag{
				Name:     "namespace",
				EnvVar:   "K3OS_OPERATOR_NAMESPACE",
				Required: true,
			},
			cli.StringFlag{
				Name:     "service-account",
				EnvVar:   "K3OS_OPERATOR_SERVICE_ACCOUNT",
				Required: true,
			},
			cli.IntFlag{
				Name:     "threads",
				EnvVar:   "K3OS_OPERATOR_THREADS",
				Required: true,
			},
		},
		Before: func(c *cli.Context) error {
			// required uid
			if os.Getuid() != 0 {
				return fmt.Errorf("must be run as root")
			}
			// required filesystem
			systemRootDir := system.RootPath()
			if inf, err := os.Stat(systemRootDir); err != nil {
				return err
			} else if !inf.IsDir() {
				return fmt.Errorf("stat %s: not a directory", systemRootDir)
			}
			return nil
		},
		Action: Run,
	}
}

// Run the `agent` sub-command
func Run(c *cli.Context) {
	logrus.Debug("K3OS::OPERATOR >>> SETUP")

	ctx := signals.SetupSignalHandler(context.Background())

	ver, err := system.GetVersion()
	if err != nil {
		logrus.Fatal(err)
	}
	if ver.Runtime != ver.Current {
		logrus.Warnf("k3os version: current(%q) != runtime(%q)", ver.Current, ver.Runtime)
	}
	logrus.Infof("k3os version: previous=%s, current=%s, runtime=%s", ver.Previous, ver.Current, ver.Runtime)

	kerndir := system.RootPath("kernel")
	if kdirinf, err := os.Stat(kerndir); err != nil {
		logrus.Warn(err)
	} else if !kdirinf.IsDir() {
		logrus.Warnf("%s is not a directory", kerndir)
	}
	kernver, err := system.GetKernelVersion()
	if err != nil {
		logrus.Warn(err)
	}
	logrus.Infof("kernel version: previous=%s, current=%s, runtime=%s", kernver.Previous, kernver.Current, kernver.Runtime)

	logrus.Debug("K3OS::OPERATOR >>> START")
	if err := controller.Start(ctx, controller.Options{
		Name:               c.String("name"),
		Namespace:          c.String("namespace"),
		ServiceAccountName: c.String("service-account"),
		Threads:            c.Int("threads"),
	}); err != nil {
		logrus.Fatalf("Error starting: %v", err)
	}

	<-ctx.Done()
}
