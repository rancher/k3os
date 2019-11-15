package upgrade

import (
	"flag"
	"runtime"
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var (
	integration bool
)

func init() {
	flag.BoolVar(&integration, "integration", false, "")
	flag.Parse()
}

func TestChannel(t *testing.T) {
	spec.Run(t, "github-releases://rancher/k3os", func(t *testing.T, when spec.G, it spec.S) {
		var (
			channel Channel
			err     error
		)
		it.Before(func() {
			if testing.Short() {
				t.Skip("-test.short=true")
			}
			if !integration {
				t.Skipf("-integration=false")
			}
			channel, err = NewChannel("github-releases://rancher/k3os")
			if err != nil {
				t.Fatal(err)
			}
		})
		when("latest", func() {
			var latest Release
			it.Before(func() {
				latest, err = channel.Latest()
				if err != nil {
					t.Fatal(err)
				}
			})
			it("version", func() {
				if latest.Name == "" {
					t.FailNow()
				}
				t.Logf("tag=%s, pre=%v", latest.Name, latest.Pre)
			})
			it("kernel-version", func() {
				if arch := runtime.GOARCH; arch == "arm" {
					t.Skipf("GOARCH=%s", arch)
				}
				ver, err := latest.KernelVersion()
				if err != nil {
					t.Error(err)
				}
				t.Logf("ver=%s", ver)
			})
		})
		when("v0.2.0", func() {
			var release Release
			it.Before(func() {
				release, err = channel.Release("v0.2.0")
				if err != nil {
					t.Fatal(err)
				}
			})
			it("version", func() {
				if release.Name == "" {
					t.FailNow()
				}
				t.Logf("tag=%s, pre=%v", release.Name, release.Pre)
			})
			it("kernel-version", func() {
				ver, err := release.KernelVersion()
				if err != nil {
					t.Error(err)
				}
				t.Logf("ver=%s", ver)
			})
		})
	}, spec.Report(report.Terminal{}))
}
