package system

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/rancher/k3os/pkg/version"
	"golang.org/x/sys/unix"
)

// Version information
type Version struct {
	Previous string
	Current  string
	Runtime  string
}

// GetVersion reads OS versioning information
func GetVersion() (Version, error) {
	ver := Version{
		Runtime: version.Version,
	}
	return ver, filesystemVersions(&ver, "k3os")
}

// GetKernelVersion returns kernel versioning information
func GetKernelVersion() (ver Version, err error) {
	if ver.Runtime, err = kernelVersion(); err != nil {
		return ver, err
	}
	return ver, filesystemVersions(&ver, "kernel")
}

func filesystemVersions(ver *Version, artifact string) error {
	current, err := os.Readlink(filepath.Join(rootDirectory, artifact, "current"))
	if err != nil {
		return err
	}
	ver.Current = filepath.Base(current)

	previous, err := os.Readlink(filepath.Join(rootDirectory, artifact, "previous"))
	if err == nil {
		ver.Previous = filepath.Base(previous)
	}

	return nil
}

func kernelVersion() (string, error) {
	utsname := unix.Utsname{}
	if err := unix.Uname(&utsname); err != nil {
		return "", err
	}
	n := bytes.IndexByte(utsname.Release[:], 0)
	return string(utsname.Release[:n]), nil
}
