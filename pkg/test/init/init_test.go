package main

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/rancher/k3os/pkg/mode"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var (
	doIntegration bool
)

func TestMain(m *testing.M) {
	flag.BoolVar(&doIntegration, "integration", false, "")
	flag.Parse()
	status := m.Run()
	fmt.Fprintf(os.Stderr, "%v\n", os.Args)
	os.Exit(status)
}

func TestPhase(t *testing.T) {

	if testing.Short() || !doIntegration {
		t.SkipNow()
	}

	spec.Run(t, "mode", func(t *testing.T, when spec.G, it spec.S) {
		it("file content should be test", func() {
			mode, err := mode.Get()
			if err != nil {
				t.Error(err)
			}
			if mode != "test" {
				t.Errorf("got %s", mode)
			}
		})
		it("envvar should be test", func() {
			mode := os.Getenv("K3OS_MODE")
			if mode != "test" {
				t.Errorf("got %s", mode)
			}
		})
		when("missing runtime file", func() {
			it.Before(func() {
				os.Rename("/run/k3os/mode", "/run/k3os/mode.tmp")
			})
			it.After(func() {
				os.Rename("/run/k3os/mode.tmp", "/run/k3os/mode")
			})
			it("should be empty", func() {
				mode, err := mode.Get()
				if err != nil {
					t.Error(err)
				}
				if mode != "" {
					t.Errorf("got %s", mode)
				}
			})
		})
	}, spec.Report(report.Terminal{}))
}
