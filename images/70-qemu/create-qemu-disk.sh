#!/bin/sh

mkdir /output

guestfish <<__EOF__
disk-create /output/k3os-qemu.qcow2 qcow2 734003200 preallocation:sparse
add-drive /output/k3os-qemu.qcow2 name:/dev/sda

disk-create "temp.img" raw 150000000 preallocation:sparse
add-drive "temp.img"

set-network true

launch

part-disk /dev/sdb mbr
mkfs ext4 /dev/sdb1
mount /dev/sdb1 /
tar-in /alpine-minirootfs-3.11.5-x86_64.tar.gz / compress:gzip

part-disk /dev/sda mbr
mkfs ext4 /dev/sda1 label:K3OS_STATE
mkdir /target
mount /dev/sda1 /target

copy-in /usr/src/qemu/. /target/
write /target/k3os/system/growpart "any any\n"

command 'apk add grub grub-efi grub-bios bash'
command "bash -c 'TARGET=/target ; DEVICE=/dev/sda ; source /target/k3os/system/k3os/current/k3os-install.sh && install_grub'"
__EOF__