package main

import (
	"os"

	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/reexec"
	"github.com/rancher/k3os/pkg/enterchroot"
	"github.com/rancher/k3os/pkg/transferroot"
	"github.com/sirupsen/logrus"
)

func main() {
	if reexec.Init() {
		return
	}

	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
	enterchroot.DebugCmdline = "k3os.debug"
	transferroot.Relocate()
	if err := mount.Mount("", "/", "none", "rw,remount"); err != nil {
		logrus.Errorf("failed to remount root as rw: %v", err)
	}
	return enterchroot.Mount("./k3os/data", os.Args, os.Stdout, os.Stderr)
}
