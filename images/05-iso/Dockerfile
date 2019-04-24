ARG REPO
ARG TAG
FROM ${REPO}/k3os-k3s:${TAG} as k3s

ARG REPO
ARG TAG
FROM ${REPO}/k3os-kernel:${TAG} as kernel

ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG} as base
ARG TAG
RUN apk add syslinux xorriso qemu-system-x86_64 libvirt qemu-img ovmf
RUN ln -s /usr/bin/qemu-system-* /usr/bin/qemu-system
RUN qemu-img create -f qcow2 /hd.img 40G

RUN mkdir -p /usr/src/iso/boot && \
    cp -rf /usr/share/syslinux /usr/src/iso/boot/isolinux
COPY syslinux/isolinux.cfg /usr/src/iso/boot/isolinux/
COPY syslinux/syslinux.cfg /usr/src/iso/boot/isolinux/
COPY syslinux/efilinux.cfg /usr/src/iso/boot/isolinux/
RUN sed -i "s/%TAG%/${TAG}/g" /usr/src/iso/boot/isolinux/isolinux.cfg \
    /usr/src/iso/boot/isolinux/syslinux.cfg \
    /usr/src/iso/boot/isolinux/efilinux.cfg

COPY --from=kernel /output/ /usr/src/kernel/
RUN cd /usr/src/kernel && \
    cp initrd /usr/src/iso/boot && \
    cp vmlinuz /usr/src/iso/boot && \
    mkdir -p /usr/src/iso/k3os/system/kernel/$(cat version) && \
    cp kernel.squashfs /usr/src/iso/k3os/system/kernel/$(cat version)

COPY --from=k3s /output/ /usr/src/k3s
RUN cd /usr/src/k3s && \
    mkdir -p /usr/src/iso/k3os/system/k3s/$(cat version) && \
    cp k3s /usr/src/iso/k3os/system/k3s/$(cat version) && \
    ln -sf $(cat version) /usr/src/iso/k3os/system/k3s/current

RUN mkdir -p /output && \
    cd /usr/src/iso && \
    xorriso -as mkisofs \
    -l -J -R -V "K3OS" \
    -no-emul-boot -boot-load-size 4 -boot-info-table \
    -b boot/isolinux/isolinux.bin -c boot/isolinux/boot.cat \
    -isohybrid-mbr boot/isolinux/isohdpfx.bin \
    -o /output/k3os.iso .

COPY run-kvm.sh /usr/bin/
CMD ["run-kvm.sh"]
