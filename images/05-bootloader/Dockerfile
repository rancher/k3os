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
RUN apk add syslinux xorriso

RUN mkdir -p /usr/src/iso/boot && \
    cp -rf /usr/share/syslinux /usr/src/iso/boot/isolinux
COPY syslinux/isolinux.cfg /usr/src/iso/boot/isolinux/
COPY syslinux/syslinux.cfg /usr/src/iso/boot/isolinux/
COPY syslinux/efilinux.cfg /usr/src/iso/boot/isolinux/
RUN sed -i "s/%VERSION%/${VERSION}/g" /usr/src/iso/boot/isolinux/isolinux.cfg \
    /usr/src/iso/boot/isolinux/syslinux.cfg \
    /usr/src/iso/boot/isolinux/efilinux.cfg

COPY --from=kernel /output/ /usr/src/kernel/
RUN cd /usr/src/kernel && \
    cp initrd /usr/src/iso/boot && \
    cp vmlinuz /usr/src/iso/boot && \
    mkdir -p /usr/src/iso/k3os/system/kernel/$(cat version) && \
    cp kernel.squashfs /usr/src/iso/k3os/system/kernel/$(cat version)

RUN mkdir -p /output && \
    cd /usr/src && \
    mv iso ${VERSION} && \
    tar cvf /output/bootloader.tar ${VERSION}
