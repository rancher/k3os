package enterchroot

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/docker/docker/pkg/mount"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func mountProc() error {
	if ok, err := mount.Mounted("/proc"); ok && err == nil {
		return nil
	}
	logrus.Debug("mkdir /proc")
	if err := os.MkdirAll("/proc", 0755); err != nil {
		return err
	}
	logrus.Debug("mount /proc")
	return mount.ForceMount("proc", "/proc", "proc", "")
}

func mountDev() error {
	if files, err := ioutil.ReadDir("/dev"); err == nil && len(files) > 2 {
		return nil
	}
	logrus.Debug("mkdir /dev")
	if err := os.MkdirAll("/dev", 0755); err != nil {
		return err
	}
	logrus.Debug("mounting /dev")
	return mount.ForceMount("none", "/dev", "devtmpfs", "")
}

func mknod(path string, mode uint32, major, minor int) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	dev := int((major << 8) | (minor & 0xff) | ((minor & 0xfff00) << 12))
	logrus.Debugf("mknod %s", path)
	return unix.Mknod(path, mode, dev)
}

func ensureloop() error {
	if err := mountProc(); err != nil {
		return errors.Wrapf(err, "failed to mount proc")
	}
	if err := mountDev(); err != nil {
		return errors.Wrapf(err, "failed to mount dev")
	}

	// ignore error
	exec.Command("modprobe", "loop").Run()

	if err := mknod("/dev/loop-control", 0700|unix.S_IFCHR, 10, 237); err != nil {
		return err
	}
	for i := 0; i < 7; i++ {
		if err := mknod(fmt.Sprintf("/dev/loop%d", i), 0700|unix.S_IFBLK, 7, i); err != nil {
			return err
		}
	}

	return nil
}
