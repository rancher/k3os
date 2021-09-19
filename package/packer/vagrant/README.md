# k3OS on Vagrant

## Quick Start

1. Build vagrant box image using [Packer](https://www.packer.io/): 

```bash
packer build .
```

2. Import vagrant box

```bash
vagrant box add --provider virtualbox k3os k3os_virtualbox.box
```

3. Run the Vagrant box:

```bash
vagrant up
```

You can then login to the box using `vagrant ssh`. See [Vagrant
Docs](https://www.vagrantup.com/docs/index.html) for more details on how
to use Vagrant

## Notes

The generated box does not have the Virtualbox Guest Additions
installed. Most of the configuration options will not work. this is
specially true for:

* `config.vm.hostname`
* `config.vm.synced_folder`
* `config.vm.network`

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
