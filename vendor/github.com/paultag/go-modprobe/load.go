package modprobe

import (
	"os"

	"golang.org/x/sys/unix"
)

// Given a short module name (such as `g_ether`), determine where the kernel
// module is located, determine any dependencies, and load all required modules.
func Load(module, params string) error {
	path, err := ResolveName(module)
	if err != nil {
		return err
	}

	order, err := Dependencies(path)
	if err != nil {
		return err
	}

	paramList := make([]string, len(order))
	paramList[len(order)-1] = params

	for i, module := range order {
		fd, err := os.Open(module)
		if err != nil {
			return err
		}
		/* not doing a defer since we're in a loop */
		param := paramList[i]
		if err := Init(fd, param); err != nil && err != unix.EEXIST {
			fd.Close()
			return err
		}
		fd.Close()
	}

	return nil
}
