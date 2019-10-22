#!/bin/bash
set -e

[ -e disk.img ] || qemu-img create -f qcow2 disk.img 8g
[ -e test.tar ] || dd if=/dev/zero of=test.tar bs=4M count=256

/usr/bin/qemu-system \
    -smp 1 \
    -m 2048 \
    -nographic \
    -display none \
    -serial mon:stdio \
    -machine accel=kvm:tcg \
    -rtc base=utc,clock=rt \
    -kernel /output/vmlinuz \
    -initrd /output/initrd \
    -nic user,restrict=off,ipv6=off,hostfwd=::2222-:22 \
    -drive if=ide,media=cdrom,index=0,file=/output/k3os.iso \
    -drive if=ide,media=cdrom,index=1,file=/output/k3os-tests.iso \
    -drive if=virtio,media=disk,file=disk.img \
    -drive if=virtio,media=disk,file=test.tar,format=raw \
    -append "loglevel=5 console=ttyS0 k3os.mode=test k3os.test.tar=/dev/vdb panic=-1" \
    -no-reboot

if tar xf test.tar --to-command 'grep -C3 FAIL' 2>/dev/null; then
    exit -1
fi
