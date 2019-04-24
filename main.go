package main

import (
	"os"

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
	return enterchroot.Mount("./k3os/data", os.Args, os.Stdout, os.Stderr)
}
