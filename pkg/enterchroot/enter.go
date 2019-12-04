package enterchroot

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/reexec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/freddierice/go-losetup.v1"
)

const (
	magic = "_SQMAGIC_"
)

var (
	symlinks     = []string{"lib", "bin", "sbin"}
	DebugCmdline = ""
)

// Enter the k3OS root
func Enter() {
	if os.Getenv("ENTER_DEBUG") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Debug("Running bootstrap")
	err := run(os.Getenv("ENTER_DATA"))
	if err != nil {
		logrus.Fatal(err)
	}
}

func isDebug() bool {
	if os.Getenv("ENTER_DEBUG") == "true" {
		return true
	}

	if DebugCmdline == "" {
		return false
	}

	bytes, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		// ignore error
		return false
	}
	for _, word := range strings.Fields(string(bytes)) {
		if word == DebugCmdline {
			return true
		}
	}

	return false
}

func Mount(dataDir string, args []string, stdout, stderr io.Writer) error {
	if err := ensureloop(); err != nil {
		return err
	}

	if isDebug() {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if logrus.GetLevel() >= logrus.DebugLevel {
		os.Setenv("ENTER_DEBUG", "true")
	}

	root, offset, err := findRoot()
	if err != nil {
		return err
	}

	os.Setenv("ENTER_DATA", dataDir)
	os.Setenv("ENTER_ROOT", root)

	logrus.Debugf("Using data [%s] root [%s]", dataDir, root)

	stat, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("failed to find %s: %v", root, err)
	}

	if !stat.IsDir() {
		logrus.Debugf("Attaching file [%s] offset [%d]", root, offset)
		dev, err := losetup.Attach(root, offset, true)
		if err != nil {
			return errors.Wrap(err, "creating loopback device")
		}
		defer dev.Detach()
		os.Setenv("ENTER_DEVICE", dev.Path())

		go func() {
			// Assume that after 3 seconds loop back device has been mounted
			time.Sleep(3 * time.Second)
			info, err := dev.GetInfo()
			if err != nil {
				return
			}

			info.Flags |= losetup.FlagsAutoClear
			err = dev.SetInfo(info)
			if err != nil {
				return
			}
		}()
	}

	logrus.Debugf("Running enter-root %v", os.Args[1:])
	if os.Getpid() == 1 {
		if err := syscall.Exec(os.Args[0], append([]string{"enter-root"}, args[1:]...), os.Environ()); err != nil {
			return errors.Wrapf(err, "failed to exec enter-root")
		}
	}
	cmd := &exec.Cmd{
		Path: os.Args[0],
		Args: append([]string{"enter-root"}, args[1:]...),
		SysProcAttr: &syscall.SysProcAttr{
			Cloneflags:   syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC,
			Unshareflags: syscall.CLONE_NEWNS,
			Pdeathsig:    syscall.SIGKILL,
		},
		Stdout: stdout,
		Stdin:  os.Stdin,
		Stderr: stderr,
		Env:    os.Environ(),
	}
	return cmd.Run()
}

func findRoot() (string, uint64, error) {
	root := os.Getenv("ENTER_ROOT")
	if root != "" {
		return root, 0, nil
	}

	for _, suffix := range []string{".root", ".squashfs"} {
		test := os.Args[0] + suffix
		if _, err := os.Stat(test); err == nil {
			return test, 0, nil
		}
	}

	return inFile()
}

func inFile() (string, uint64, error) {
	f, err := os.Open(reexec.Self())
	if err != nil {
		return "", 0, err
	}
	defer f.Close()

	buf := make([]byte, 8192)
	test := []byte(strings.ToLower(magic))
	testLength := len(test)
	offset := uint64(0)
	found := 0

	for {
		n, err := f.Read(buf)
		if err == io.EOF && n == 0 {
			break
		} else if err != nil {
			return "", 0, err
		}

		for _, b := range buf[:n] {
			if b == test[found] {
				found++
				if found == testLength {
					return reexec.Self(), offset + 1, nil
				}
			} else {
				found = 0
			}
			offset++
		}
	}

	return "", 0, fmt.Errorf("failed to find image in file %s", os.Args[0])
}

func run(data string) error {
	mounted, err := mount.Mounted(data)
	if err != nil {
		return errors.Wrapf(err, "checking %s mounted", data)
	}

	if !mounted {
		if err := os.MkdirAll(data, 0755); err != nil {
			return errors.Wrapf(err, "mkdir %s", data)
		}
		if err := mount.Mount(data, data, "none", "rbind"); err != nil {
			return errors.Wrapf(err, "remounting data %s", data)
		}
	}

	root := os.Getenv("ENTER_ROOT")
	device := os.Getenv("ENTER_DEVICE")

	logrus.Debugf("Using root %s %s", root, device)

	usr := filepath.Join(data, "usr")
	dotRoot := filepath.Join(data, ".base")

	for _, d := range []string{usr, dotRoot} {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("failed to make dir %s: %v", data, err)
		}
	}

	if device == "" {
		logrus.Debugf("Bind mounting %s to %s", root, usr)
		if err := mount.Mount(root, usr, "none", "bind"); err != nil {
			return fmt.Errorf("failed to bind mount")
		}
	} else {
		logrus.Debugf("Mounting squashfs %s to %s", device, usr)
		squashErr := checkSquashfs()
		err = mount.Mount(device, usr, "squashfs", "ro")
		if err != nil {
			err = errors.Wrap(err, "mounting squashfs")
			if squashErr != nil {
				err = errors.Wrap(err, squashErr.Error())
			}
			return err
		}
	}

	if err := os.Chdir(data); err != nil {
		return err
	}

	for _, p := range symlinks {
		if _, err := os.Lstat(p); os.IsNotExist(err) {
			if err := os.Symlink(filepath.Join("usr", p), p); err != nil {
				return errors.Wrapf(err, "failed to symlink %s", p)
			}
		}
	}

	logrus.Debugf("pivoting to . .base")
	if err := syscall.PivotRoot(".", ".base"); err != nil {
		return errors.Wrap(err, "pivot_root failed")
	}

	if err := mount.ForceMount("", ".", "none", "rprivate"); err != nil {
		return errors.Wrapf(err, "making . private %s", data)
	}

	if err := syscall.Chroot("/"); err != nil {
		return err
	}

	if err := os.Chdir("/"); err != nil {
		return err
	}

	if _, err := os.Stat("/usr/init"); err != nil {
		return errors.Wrap(err, "failed to find /usr/init")
	}

	os.Unsetenv("ENTER_ROOT")
	os.Unsetenv("ENTER_DATA")
	os.Unsetenv("ENTER_DEVICE")
	return syscall.Exec("/usr/init", os.Args, os.Environ())
}

func checkSquashfs() error {
	if !inProcFS() {
		exec.Command("modprobe", "squashfs").Run()
	}

	if !inProcFS() {
		return errors.New("This kernel does not support squashfs, please enable. " +
			"On Fedora you may need to run \"dnf install kernel-modules-$(uname -r)\"")
	}

	return nil
}

func inProcFS() bool {
	bytes, err := ioutil.ReadFile("/proc/filesystems")
	if err != nil {
		logrus.Errorf("Failed to read /proc/filesystems: %v", err)
		return false
	}
	return strings.Contains(string(bytes), "squashfs")
}
