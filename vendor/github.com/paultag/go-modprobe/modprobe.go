package modprobe

import (
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Given an open .ko file's os.File (created with os.Open or similar), attempt
// to load that kernel module into the running kernel. This may error out
// for a number of reasons, such as no permission (either setcap CAP_SYS_MODULE
// or run as root), the .ko being for the wrong kernel, or the file not being
// a module at all.
//
// Any arguments to the module may be passed through `params`, such as
// `file=/root/data/backing_file`.
func Init(file *os.File, params string) error {
	_p0, err := unix.BytePtrFromString(params)
	if err != nil {
		return err
	}

	if _, _, err := unix.Syscall(unix.SYS_FINIT_MODULE, file.Fd(), uintptr(unsafe.Pointer(_p0)), 0); err != 0 {
		return err
	}
	return nil
}

// Unload a loaded kernel module. If no such module is loaded, or if the module
// can not be unloaded, this function will return an error.
func Remove(name string) error {
	moduleName, err := unix.BytePtrFromString(name)
	if err != nil {
		return err
	}

	if _, _, err := unix.Syscall(unix.SYS_DELETE_MODULE, uintptr(unsafe.Pointer(moduleName)), 0, 0); err != 0 {
		return err
	}
	return nil
}
