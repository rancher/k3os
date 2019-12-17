package upgrade

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/docker/docker/pkg/mount"
	"github.com/rancher/k3os/pkg/system"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/sys/unix"
)

// Command is the `upgrade` sub-command, it performs upgrades to k3OS.
func Command() cli.Command {
	return cli.Command{
		Name:  "upgrade",
		Usage: "perform upgrades",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:   "kernel",
				Usage:  "upgrade the kernel",
				EnvVar: "K3OS_UPGRADE_KERNEL",
			},
			cli.BoolFlag{
				Name:   "rootfs",
				Usage:  "upgrade k3os+k3s",
				EnvVar: "K3OS_UPGRADE_ROOTFS",
			},
			cli.BoolFlag{
				Name:   "remount",
				Usage:  "pre-upgrade remount?",
				EnvVar: "K3OS_UPGRADE_REMOUNT",
			},
			cli.BoolFlag{
				Name:   "sync",
				Usage:  "post-upgrade sync?",
				EnvVar: "K3OS_UPGRADE_SYNC",
			},
			cli.StringFlag{
				Name:     "source",
				EnvVar:   "K3OS_UPGRADE_SOURCE",
				Value:    system.RootPath(),
				Required: true,
			},
			cli.StringFlag{
				Name:     "destination",
				EnvVar:   "K3OS_UPGRADE_DESTINATION",
				Value:    system.RootPath(),
				Required: true,
			},
			cli.StringFlag{
				Name:   "lock-file",
				Value:  system.StatePath("upgrade.lock"),
				Hidden: true,
			},
		},
		Action: Run,
		Before: func(c *cli.Context) error {
			if dst, src := c.String("destination"), c.String("source"); src == dst {
				cli.ShowSubcommandHelp(c)
				logrus.Errorf("the `source` cannot be the `destination`: %s", src)
				os.Exit(1)
			}
			if !c.Bool("rootfs") && !c.Bool("kernel") {
				cli.ShowSubcommandHelp(c)
				logrus.Error("at least one of `rootfs` or `kernel` must be true")
				os.Exit(1)
			}
			return nil
		},
		After: func(c *cli.Context) error {
			if c.Bool("sync") {
				unix.Sync()
			}
			return nil
		},
	}
}

// Run the `upgrade` sub-command
func Run(c *cli.Context) {
	lockFile := c.String("lock-file")
	//defer os.Remove(lockFile)

	// establish the lock
	lf, err := os.OpenFile(lockFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	defer lf.Close()
	if err = unix.Flock(int(lf.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		logrus.Fatal(err)
	}
	defer unix.Flock(int(lf.Fd()), unix.LOCK_UN)

	sourceDir := c.String("source")
	if err := validateSystemRoot(sourceDir); err != nil {
		logrus.Fatal(err)
	}
	destinationDir := c.String("destination")
	if err := validateSystemRoot(destinationDir); err != nil {
		logrus.Fatal(err)
	}

	if c.Bool("remount") {
		if err := mount.Mount("", destinationDir, "none", "remount,rw"); err != nil {
			logrus.Fatal(err)
		}
	}

	if c.Bool("kernel") {
		if err := copyArtifact(sourceDir, destinationDir, "kernel"); err != nil {
			logrus.Fatal(err)
		}
	}

	if c.Bool("rootfs") {
		if err := copyArtifact(sourceDir, destinationDir, "k3s"); err != nil {
			logrus.Fatal(err)
		}
		if err := copyArtifact(sourceDir, destinationDir, "k3os"); err != nil {
			logrus.Fatal(err)
		}
	}
}

func copyArtifact(src, dst, artifact string) error {
	rsync := exec.Command("rsync",
		"-av", fmt.Sprintf("%s/%s/", src, artifact), fmt.Sprintf("%s/%s/", dst, artifact),
	)
	rsync.Stderr = os.Stderr
	rsync.Stdout = os.Stdout
	return rsync.Run()
}

func validateSystemRoot(systemDir string) error {
	inf, err := os.Stat(systemDir)
	if err != nil {
		return err
	}
	if !inf.IsDir() {
		return fmt.Errorf("stat %s: not a directory", systemDir)
	}
	return nil
}
