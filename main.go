package main

import (
	"os"

	"github.com/niusmallnan/k3os/cmd/control"

	"github.com/moby/moby/pkg/reexec"
	"github.com/sirupsen/logrus"
)

var entryPoints = map[string]func(){
	"k3os-sysinit": control.SysInitMain,
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
