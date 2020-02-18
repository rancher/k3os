# go-modprobe

Load an unload Linux kernel modules using the Linux module syscalls.

This package is Linux specific. Loading a module uses the `finit` variant,
which allows loading of modules by a file descriptor, rather than having to
load an ELF into the process memory before loading.

The ability to load and unload modules is dependent on either the `CAP_SYS_MODULE`
capability, or running as root. Care should be taken to understand what security
implications this has on processes that use this library.

## Setting the capability on a binary using this package

```
$ sudo setcap cap_sys_module+ep /path/to/binary
```
