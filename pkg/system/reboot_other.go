// +build !linux

package system

// reboot exists to make my ide less red
func reboot() {
	panic("REBOOT")
}
