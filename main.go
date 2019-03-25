package main

import (
	"os"

	"github.com/niusmallnan/k3os/cmd/control"

	"github.com/sirupsen/logrus"
)

func init() {
	// TODO: need to remove just for develop
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	control.Main()
}
