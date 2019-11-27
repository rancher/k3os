// +build linux

package system

import (
	"golang.org/x/sys/unix"
)

// reboot exists to make my ide less red
func reboot() {
	unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART)
}
