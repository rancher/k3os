ARG REPO
ARG TAG
FROM ${REPO}/k3os-iso:${TAG} as iso

ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG} as base
ARG VERSION
ARG ARCH

RUN apk add xorriso grub grub-efi mtools libvirt qemu-img dnsmasq
RUN if [ "$ARCH" == "amd64" ]; then \
        apk add qemu-system-x86_64 grub-bios ovmf \
    ;elif [ "$ARCH" == "arm64" ]; then \
        apk add qemu-system-aarch64 \
    ;fi
RUN ln -s /usr/bin/qemu-system-* /usr/bin/qemu-system
RUN qemu-img create -f qcow2 /hd.img 40G

COPY run-kvm.sh /usr/bin/
COPY qemu-ifup /etc/
COPY qemu-ifdown /etc/
RUN chmod 755 /etc/qemu-ifup
RUN chmod 755 /etc/qemu-ifdown

RUN mkdir -p /output
COPY --from=iso /output/k3os.iso /output/

CMD ["run-kvm.sh"]
