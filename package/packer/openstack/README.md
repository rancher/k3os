# k3os/package/packer/openstack

This is an Example for building a K3OS OpenStack Image using packer.  

## Quick start

## Notes

* I had to build packer from source for aarch64 due to the following issue with packer: [8258](https://github.com/hashicorp/packer/issues/8258) 


### Empty log output for linux instances


```

console=tty0 console=ttyS0,115200n8

```
## References

* [https://github.com/hashicorp/packer/issues/8258](https://github.com/hashicorp/packer/issues/8258)
