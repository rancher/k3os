this fixture is missing both the `current` and `previous` symlinks which should cause both
`system.GetVersion()` and `system.GetKernelVersion()` to return an error.