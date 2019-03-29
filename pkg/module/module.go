package module

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/niusmallnan/k3os/config"

	"github.com/sirupsen/logrus"
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
	modules := cfg.K3OS.Defaults.Modules
	additional := cfg.K3OS.Modules
	modules = append(modules, additional...)
	for _, m := range modules {
		if loaded[m] {
			continue
		}
		params := strings.SplitN(m, " ", -1)
		logrus.Debugf("module %s with parameters [%s] is loading", m, params)
		cmd := exec.Command("modprobe", params...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("could not load module %s with parameters [%s], err %v", m, params, err)
		}
		logrus.Debugf("module %s is loaded", m)
	}
	return sc.Err()
}
