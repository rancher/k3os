package system

import "time"

// Reboot the system after the specified delay
func Reboot(delay time.Duration) {
	time.Sleep(delay)
	reboot()
}
