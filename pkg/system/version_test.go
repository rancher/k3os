package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && filepath.Dir(path) == "testdata" {
			rootDirectory = path
			base := filepath.Base(path)
			t.Run(base, func(t *testing.T) {
				ver, err := GetVersion()
				if strings.HasPrefix(base, "error") {
					if err == nil {
						t.Fatal("expecting error")
					}
					t.Log(err)
				} else if strings.HasPrefix(base, "noerr") {
					if err != nil {
						t.Fatal(err)
					}
					t.Logf("previous=%s, current=%s, running=%s", ver.Previous, ver.Current, ver.Runtime)
					if strings.HasSuffix(base, "missing-previous") && ver.Previous != "" {
						t.Fatal("not missing previous")
					} else if strings.HasSuffix(base, "with-previous") && ver.Previous == "" {
						t.Fatal("not with previous")
					} else if ver.Current == "" {
						t.Fatal("not with current")
					} else if ver.Runtime == "" {
						t.Fatal("not with running")
					}
				}
			})
		}
		return nil
	})
}

func TestGetKernelVersion(t *testing.T) {
	filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && filepath.Dir(path) == "testdata" {
			rootDirectory = path
			base := filepath.Base(path)
			t.Run(base, func(t *testing.T) {
				ver, err := GetKernelVersion()
				if strings.HasPrefix(base, "error") {
					if err == nil {
						t.Fatal("expecting error")
					}
					t.Log(err)
				} else if strings.HasPrefix(base, "noerr") {
					if err != nil {
						t.Fatal(err)
					}
					t.Logf("previous=%s, current=%s, runtime=%s", ver.Previous, ver.Current, ver.Runtime)
					if strings.HasSuffix(base, "missing-previous") && ver.Previous != "" {
						t.Fatal("not missing previous")
					} else if strings.HasSuffix(base, "with-previous") && ver.Previous == "" {
						t.Fatal("not with previous")
					} else if ver.Current == "" {
						t.Fatal("not with current")
					} else if ver.Runtime == "" {
						t.Fatal("not with running")
					}
				}
			})
		}
		return nil
	})
}
