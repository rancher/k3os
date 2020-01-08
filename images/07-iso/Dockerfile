ARG REPO
ARG TAG
FROM ${REPO}/k3os-package:${TAG} as package

ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG} as base
ARG VERSION
ARG ARCH
RUN apk add xorriso grub grub-efi mtools libvirt qemu-img
RUN if [ "$ARCH" == "amd64" ]; then \
        apk add qemu-system-x86_64 grub-bios ovmf \
    ;elif [ "$ARCH" == "arm64" ]; then \
        apk add qemu-system-aarch64 \
    ;fi
RUN ln -s /usr/bin/qemu-system-* /usr/bin/qemu-system
RUN qemu-img create -f qcow2 /hd.img 40G
COPY run-kvm.sh /usr/bin/
COPY grub.cfg /usr/src/iso/boot/grub/grub.cfg

COPY --from=package /output/ /usr/src/iso/

COPY wrapper /usr/bin/
COPY config.yaml /usr/src/iso/k3os/system/
RUN mkdir -p /output && \
    cd /usr/src/iso && \
    grub-mkrescue --xorriso=/usr/bin/wrapper -o /output/k3os.iso . -V K3OS && \
    [ -e /output/k3os.iso ] # grub-mkrescue doesn't exit non-zero on failure

CMD ["run-kvm.sh"]
