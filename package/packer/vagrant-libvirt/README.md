# k3OS on Vagrant using the libvirt provider

## Quick Start

1. Build vagrant box image for libvirt using [Packer](https://www.packer.io/): 

```bash
packer build .
```

2. Import vagrant box

```bash
vagrant box add --provider libvirt k3os k3os_libvirt.box
```

3. Run the Vagrant box:

```bash
vagrant up
```

You can then login to the box using `vagrant ssh`. See [Vagrant
Docs](https://www.vagrantup.com/docs/index.html) for more details on how
to use Vagrant

## Notes

The shell provisioner is working but requires some tweaking. The
provisioning shell script will be put into `/tmp/vagrant-shell`. The
file can be written but can not be executed. To mitigate this
limitation, one have to set the `upload_path` option.

```
config.vm.provision 'shell',
  upload_path: '/home/rancher/vagrant-shell',
  inline: <<-SHELL
mkdir -p /mnt
mount /dev/sda1 /mnt
cat <<EOF > /mnt/k3os/system/config.yaml
k3os:
  token: EbvX0V38syjPQBZJ71tb9EIHbyL5mISBqDSTa2aJt7LSCF1JEW
  password: rancher
EOF
reboot
SHELL
```

The above example also shows how the k3OS config can be changed. When
you do so, you have to set the password to `rancher`. If this is not the
case, vagrant will not be able to login.
