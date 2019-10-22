#!/bin/bash
set -e

[ -e disk.img ] || qemu-img create -f qcow2 disk.img 40G

exec /usr/bin/qemu-system \
    -smp 1 \
    -m 2048 \
    -nographic \
    -display none \
    -serial mon:stdio \
    -machine accel=kvm:tcg \
    -rtc base=utc,clock=rt \
    -nic user,restrict=off,ipv6=off,hostfwd=::2222-:22 \
    -drive if=ide,media=cdrom,index=0,file=/output/k3os.iso \
    -drive if=virtio,media=disk,file=disk.img
