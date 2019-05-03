ARG REPO
ARG TAG
FROM ${REPO}/k3os-tar:${TAG} as tar

ARG REPO
ARG TAG
FROM ${REPO}/k3os-kernel:${TAG} as kernel

ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG} as base
ARG VERSION
RUN apk add xorriso grub grub-efi grub-bios
RUN apk add qemu-system-x86_64 libvirt qemu-img ovmf
RUN ln -s /usr/bin/qemu-system-* /usr/bin/qemu-system
RUN qemu-img create -f qcow2 /hd.img 40G
COPY run-kvm.sh /usr/bin/
COPY grub.cfg /usr/src/iso/boot/grub/grub.cfg

COPY --from=kernel /output/ /usr/src/kernel/
RUN cd /usr/src/kernel && \
    mkdir -p /usr/src/iso/k3os/system/kernel/$(cat version) && \
    cp initrd kernel.squashfs /usr/src/iso/k3os/system/kernel/$(cat version) && \
    ln -s $(cat version) /usr/src/iso/k3os/system/kernel/current

COPY --from=tar /output/userspace.tar /usr/src/tars/
RUN tar xvf /usr/src/tars/userspace.tar --strip-components=1 -C /usr/src/iso

RUN mkdir -p /output && \
    cd /usr/src/iso && \
    grub-mkrescue -o /output/k3os.iso . -V K3OS && \
    [ -e /output/k3os.iso ] # grub-mkrescue doesn't exit non-zero on failure

CMD ["run-kvm.sh"]
