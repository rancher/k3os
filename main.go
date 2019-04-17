package main

import (
	"os"

	"github.com/rancher/k3os/cmd/control"
	"github.com/rancher/k3os/cmd/netinit"
	"github.com/rancher/k3os/cmd/sysinit"
	"github.com/rancher/k3os/cmd/ttyinit"

	"github.com/docker/docker/pkg/reexec"
	"github.com/sirupsen/logrus"
)

var entryPoints = map[string]func(){
	"/usr/lib/k3os/k3os-sysinit": sysinit.Main,
	"/usr/lib/k3os/k3os-netinit": netinit.Main,
	"/usr/lib/k3os/k3os-ttyinit": ttyinit.Main,
}

func init() {
	// TODO: need to remove just for develop
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	for n, f := range entryPoints {
		reexec.Register(n, f)
	}
	if !reexec.Init() {
		control.Main()
	}
}
