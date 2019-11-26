package main

// Copyright 2019 Rancher Labs, Inc.
// SPDX-License-Identifier: Apache-2.0

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/reexec"
	"github.com/rancher/k3os/pkg/cli/app"
	"github.com/rancher/k3os/pkg/enterchroot"
	"github.com/rancher/k3os/pkg/transferroot"
	"github.com/sirupsen/logrus"
)

func main() {
	reexec.Register("/init", initrd)      // mode=live
	reexec.Register("/sbin/init", initrd) // mode=local
	reexec.Register("enter-root", enterchroot.Enter)

	if !reexec.Init() {
		app := app.New()
		args := []string{app.Name}
		path := filepath.Base(os.Args[0])
		if path != app.Name && app.Command(path) != nil {
			args = append(args, path)
		}
		args = append(args, os.Args[1:]...)
		// this will bomb if the app has any non-defaulted, required flags
		err := app.Run(args)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func initrd() {
	enterchroot.DebugCmdline = "k3os.debug"
	transferroot.Relocate()
	if err := mount.Mount("", "/", "none", "rw,remount"); err != nil {
		logrus.Errorf("failed to remount root as rw: %v", err)
	}
	if err := enterchroot.Mount("./k3os/data", os.Args, os.Stdout, os.Stderr); err != nil {
		logrus.Fatalf("failed to enter root: %v", err)
	}
}
