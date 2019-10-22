ARG REPO
ARG TAG
FROM ${REPO}/k3os-kernel:${TAG} as kernel

ARG REPO
ARG TAG
FROM ${REPO}/k3os-iso:${TAG} as iso
ARG VERSION
ARG ARCH

RUN set -x \
 && apk add \
    libvirt \
    qemu-img
RUN [ "$ARCH" == "amd64" ] \
 && set -x \
 && apk add \
    grub-bios \
    ovmf \
    qemu-system-x86_64 \
 || true
RUN [ "$ARCH" == "arm64" ] \
 && set -x \
 && apk add \
    qemu-system-aarch64 \
 || true
RUN ln -s /usr/bin/qemu-system-* /usr/bin/qemu-system

COPY --from=kernel /output/vmlinuz /output/
COPY --from=kernel /output/initrd /output/
COPY --from=kernel /output/kernel.squashfs /output/
COPY *.sh /usr/local/bin/

WORKDIR /output
