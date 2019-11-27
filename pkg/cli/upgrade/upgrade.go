package upgrade

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/pkg/mount"
	"github.com/rancher/k3os/pkg/system"
	"github.com/rancher/k3os/pkg/upgrade"
	"github.com/rancher/k3os/pkg/util"
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
			cli.GenericFlag{
				Name:   "channel",
				Usage:  "updates/releases uri",
				EnvVar: "K3OS_UPGRADE_CHANNEL",
				Value: &Channel{
					uri: "https://github.com/rancher/k3os/releases/latest",
				},
			},
			cli.StringFlag{
				Name:   "version",
				Usage:  "release or tag",
				EnvVar: "K3OS_UPGRADE_VERSION",
				Value:  "latest",
			},
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
			cli.BoolFlag{
				Name:   "reboot",
				Usage:  "post-upgrade reboot?",
				EnvVar: "K3OS_UPGRADE_REBOOT",
			},
			cli.StringFlag{
				Name:   "system-dir",
				Value:  system.RootPath(),
				Hidden: true,
			},
			cli.StringFlag{
				Name:   "lock-file",
				Value:  system.StatePath("upgrade.lock"),
				Hidden: true,
			},
		},
		Action: Run,
		Before: func(c *cli.Context) error {
			chn := c.Generic("channel").(*Channel)
			if chn != nil && chn.upchan == nil {
				if err := chn.Set(chn.uri); err != nil {
					return err
				}
			}
			if !c.Bool("rootfs") && !c.Bool("kernel") {
				cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
			}
			return nil
		},
		After: func(c *cli.Context) error {
			if c.Bool("sync") {
				unix.Sync()
			}
			if c.Bool("reboot") {
				system.Reboot(0 * time.Second)
			}
			return nil
		},
	}
}

// Run the `upgrade` sub-command
func Run(c *cli.Context) {
	lockFile := c.String("lock-file")
	version := c.String("version")

	// establish the lock
	lf, err := os.OpenFile(lockFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	defer lf.Close()
	if err = unix.Flock(int(lf.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		logrus.Fatal(err)
	}
	if _, err := lf.WriteString(version); err != nil {
		logrus.Warnf("unable to write version to lock-file: %v", err)
	}
	defer unix.Flock(int(lf.Fd()), unix.LOCK_UN)

	// establish that we can write to the system-dir
	systemDir := c.String("system-dir")
	if inf, err := os.Stat(systemDir); err != nil {
		logrus.Fatal(err)
	} else if !inf.IsDir() {
		logrus.Fatalf("stat %s: not a directory", systemDir)
	}
	if c.Bool("remount") {
		if err := mount.Mount("", systemDir, "none", "remount,rw"); err != nil {
			logrus.Fatal(err)
		}
	}

	channel := c.Generic("channel").(*Channel)
	release, err := channel.Release(version)
	if err != nil {
		logrus.Fatal(err)
	}

	if c.Bool("kernel") {
		if err := downloadKernel(systemDir, &release); err != nil {
			logrus.Fatal(err)
		}
	}

	if c.Bool("rootfs") {
		if err := downloadRootfs(systemDir, &release); err != nil {
			logrus.Fatal(err)
		}
	}

}

func downloadKernel(dir string, rel *upgrade.Release) error {
	dir = filepath.Join(dir, "kernel")
	if inf, err := os.Stat(dir); err != nil {
		return err
	} else if !inf.IsDir() {
		return fmt.Errorf("stat %s: not a directory", dir)
	}

	krdn := "k3os-initrd-" + runtime.GOARCH
	krda := rel.Asset(krdn)
	if krda == nil {
		return fmt.Errorf("asset not found: %s", krdn)
	}
	ksqn := "k3os-kernel-" + runtime.GOARCH + ".squashfs"
	ksqa := rel.Asset(ksqn)
	if ksqa == nil {
		return fmt.Errorf("asset not found: %s", ksqn)
	}

	ver, err := rel.KernelVersion()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, ver)
	if err = os.MkdirAll(path, 0555); err != nil {
		return err
	}

	err = util.Stream(krda.URL, writeFile(filepath.Join(path, "initrd"), 0644))
	if err != nil {
		return err
	}
	err = util.Stream(ksqa.URL, writeFile(filepath.Join(path, "kernel.squashfs"), 0644))
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cwd)

	if err = os.Chdir(filepath.Dir(path)); err != nil {
		return err
	}
	if err = os.Rename("current", "previous"); err != nil {
		logrus.Warn(err)
	}
	if err = os.Symlink(ver, "current"); err != nil {
		return err
	}

	return nil
}

func downloadRootfs(dir string, rel *upgrade.Release) error {
	rootfsName := "k3os-rootfs-" + runtime.GOARCH + ".tar.gz"
	rootfsAsset := rel.Asset(rootfsName)
	if rootfsAsset == nil {
		return fmt.Errorf("asset not found: %s", rootfsName)
	}
	return util.Stream(rootfsAsset.URL, func(body io.Reader) error {
		gzr, err := gzip.NewReader(body)
		if err != nil {
			return err
		}
		defer gzr.Close()
		tr := tar.NewReader(gzr)
		for {
			header, err := tr.Next()
			switch err {
			case nil:
				// process the header (see below)
			case io.EOF:
				return nil
			default:
				return err
			}
			if strings.HasPrefix(header.Name, filepath.Join(rel.Name, system.DefaultRootDir)) {
				path := filepath.Join(dir, filepath.Join(strings.Split(header.Name, "/")[3:]...))
				switch header.Typeflag {
				case tar.TypeDir:
					logrus.Debugf("# %s", path)
					if err = os.MkdirAll(path, 0555); err != nil {
						return err
					}
				case tar.TypeSymlink:
					if dir, base := filepath.Split(path); base == "current" {
						defer func() {
							logrus.Debugf("# %s", path)
							if err = os.Rename(path, filepath.Join(dir, "previous")); err != nil {
								logrus.Warn(err)
							}
							if err = os.Symlink(header.Linkname, path); err != nil {
								return
							}
						}()
					} else {
						logrus.Debugf("# %s", path)
						if err = os.Symlink(header.Linkname, path); err != nil {
							return err
						}
					}
				case tar.TypeReg:
					writeFrom := writeFile(path, os.FileMode(header.Mode))
					if err = writeFrom(tr); err != nil {
						return err
					}
				}
			}
		}
	})
}

func writeFile(name string, mode os.FileMode) func(body io.Reader) error {
	return func(body io.Reader) error {
		file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, mode)
		if err != nil {
			return err
		}
		defer file.Close()
		logrus.Debugf("# > %s", file.Name())
		_, err = io.Copy(file, body)
		return err
	}
}
