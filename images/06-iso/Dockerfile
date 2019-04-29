ARG REPO
ARG TAG
FROM ${REPO}/k3os-tar:${TAG} as tar

ARG REPO
ARG TAG
FROM ${REPO}/k3os-kernel:${TAG} as kernel

ARG REPO
ARG TAG
FROM ${REPO}/k3os-bootloader:${TAG} as bootloader

ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG} as base
ARG VERSION
RUN apk add syslinux xorriso

COPY --from=bootloader /output/bootloader.tar /usr/src/tars/
COPY --from=tar /output/userspace.tar /usr/src/tars/
RUN mkdir -p /usr/src/iso && \
    tar xvf /usr/src/tars/bootloader.tar --strip-component=1 -C /usr/src/iso && \
    tar xvf /usr/src/tars/userspace.tar --strip-components=1 --exclude=k3os/system/k3os -C /usr/src/iso

RUN mkdir -p /output && \
    cd /usr/src/iso && \
    curl -o boot/isolinux/efiboot.img -fL http://slackware.mirrors.tds.net/pub/slackware/slackware64-current/isolinux/efiboot.img && \
    xorriso -as mkisofs \
    -V "K3OS" \
    -o /output/k3os.iso \
    -isohybrid-mbr boot/isolinux/isohdpfx.bin \
    -eltorito-alt-boot \
    -eltorito-platform 0xEF \
    -e boot/isolinux/efiboot.img \
    -no-emul-boot \
    -isohybrid-gpt-basdat \
    -o /output/k3os.iso .
