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
func GetVersion(prefix ...string) (ver Version, err error) {
	if len(prefix) == 0 {
		ver.Runtime = version.Version
	}
	return ver, filesystemVersions(&ver, filepath.Join(prefix...), "k3os")
}

// GetKernelVersion returns kernel versioning information
func GetKernelVersion(prefix ...string) (ver Version, err error) {
	if len(prefix) == 0 {
		if ver.Runtime, err = unameRelease(); err != nil {
			return ver, err
		}
	}
	return ver, filesystemVersions(&ver, filepath.Join(prefix...), "kernel")
}

func filesystemVersions(ver *Version, prefix, artifact string) error {
	currentPath := filepath.Join(filepath.Join(prefix), filepath.Join(rootDirectory, artifact, "current"))
	current, err := os.Readlink(currentPath)
	if err != nil {
		return err
	}
	ver.Current = filepath.Base(current)

	previousPath := filepath.Join(filepath.Join(prefix), filepath.Join(rootDirectory, artifact, "previous"))
	previous, err := os.Readlink(previousPath)
	if err == nil {
		ver.Previous = filepath.Base(previous)
	}

	return nil
}

func unameRelease() (string, error) {
	utsname := unix.Utsname{}
	if err := unix.Uname(&utsname); err != nil {
		return "", err
	}
	n := bytes.IndexByte(utsname.Release[:], 0)
	return string(utsname.Release[:n]), nil
}
