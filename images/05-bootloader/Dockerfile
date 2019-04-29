ARG REPO
ARG TAG
FROM ${REPO}/k3os-k3s:${TAG} as k3s

ARG REPO
ARG TAG
FROM ${REPO}/k3os-kernel:${TAG} as kernel

ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG} as base
ARG VERSION
RUN apk add syslinux xorriso grub grub-efi

RUN mkdir -p /usr/src/iso/boot && \
    cp -rf /usr/share/syslinux /usr/src/iso/boot/isolinux
COPY syslinux/isolinux.cfg /usr/src/iso/boot/isolinux/
COPY syslinux/syslinux.cfg /usr/src/iso/boot/isolinux/
COPY syslinux/efilinux.cfg /usr/src/iso/boot/isolinux/
COPY syslinux/grub.cfg /usr/src/iso/boot/grub/
RUN sed -i "s/%VERSION%/${VERSION}/g" /usr/src/iso/boot/isolinux/isolinux.cfg \
    /usr/src/iso/boot/isolinux/syslinux.cfg \
    /usr/src/iso/boot/isolinux/efilinux.cfg

COPY --from=kernel /output/ /usr/src/kernel/
RUN cd /usr/src/kernel && \
    cp initrd /usr/src/iso/boot && \
    cp vmlinuz /usr/src/iso/boot && \
    mkdir -p /usr/src/iso/k3os/system/kernel/$(cat version) && \
    cp kernel.squashfs /usr/src/iso/k3os/system/kernel/$(cat version)

RUN mkdir -p /tmp/efi/boot && \
    grub-mkimage \
        -C xz \
        -O x86_64-efi \
        -p /boot/grub \
        -o /tmp/efi/boot/bootx64.efi \
        boot linux search normal configfile \
        part_gpt btrfs fat iso9660 loopback \
        test keystatus gfxmenu regexp probe \
        efi_gop efi_uga all_video gfxterm font \
        echo read ls cat png jpeg halt reboot && \
    truncate -s 4M /usr/src/iso/boot/isolinux/efiboot.img && \
    mkfs.vfat /usr/src/iso/boot/isolinux/efiboot.img && \
    mcopy -i /usr/src/iso/boot/isolinux/efiboot.img -s /tmp/efi ::

RUN mkdir -p /output && \
    cd /usr/src && \
    mv iso ${VERSION} && \
    tar cvf /output/bootloader.tar ${VERSION}

