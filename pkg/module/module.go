package module

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rancher/k3os/pkg/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

const (
	procModulesFile = "/proc/modules"
)

func LoadModules(cfg *config.CloudConfig) error {
	loaded := map[string]bool{}
	f, err := os.Open(procModulesFile)
	if err != nil {
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		loaded[strings.SplitN(sc.Text(), " ", 2)[0]] = true
	}
	modules := cfg.K3OS.Modules
	for _, m := range modules {
		if loaded[m] {
			continue
		}
		params := strings.SplitN(m, " ", -1)
		logrus.Debugf("module %s with parameters [%s] is loading", m, params)
		if err := modprobe(params[0], params[1:]); err != nil {
			return fmt.Errorf("could not load module %s with parameters [%s], err %v", m, params, err)
		}
		logrus.Debugf("module %s is loaded", m)
	}
	return sc.Err()
}

func modprobe(module string, params []string) error {
	uname := unix.Utsname{}
	if err := unix.Uname(&uname); err != nil {
		return fmt.Errorf("unable to determine uname, err %v", err)
	}
	i := 0
	for ; uname.Release[i] != 0; i++ {
	}
	pth := fmt.Sprintf("/lib/modules/%s/**/%s.ko", uname.Release[:i], module)
	files, err := filepath.Glob(pth)
	if err != nil {
		return fmt.Errorf("unable to search for module, err %v", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("module not found")
	}
	file, err := os.Open(files[0])
	if err != nil {
		return fmt.Errorf("could not open module file %s, err %v", files[0], err)
	}
	if err := unix.FinitModule(int(file.Fd()), strings.Join(params, " "), 0); err != nil {
		return fmt.Errorf("unable to load module, err %v", err)
	}
	return nil
}
