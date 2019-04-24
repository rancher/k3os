ARG REPO
ARG TAG
FROM ${REPO}/k3os-bin:${TAG} as bin

FROM ${REPO}/k3os-kernel-stage1:${TAG} as kernel

FROM ${REPO}/k3os-base:${TAG}
RUN apk add squashfs-tools
COPY --from=kernel /output/ /usr/src/kernel/

RUN mkdir -p /usr/src/initrd/lib && \
    cd /usr/src/kernel && \
    tar cf - -T initrd-modules -T initrd-firmware | tar xf - -C /usr/src/initrd/ && \
    depmod -b /usr/src/initrd $(cat /usr/src/kernel/version)

RUN mkdir -p /output && \
    cd /usr/src/kernel && \
    depmod -b . $(cat /usr/src/kernel/version) && \
    mksquashfs . /output/kernel.squashfs

RUN cp /usr/src/kernel/version /output/ && \
    cp /usr/src/kernel/vmlinuz /output/

COPY --from=bin /output/k3os /usr/src/initrd/init
RUN cd /usr/src/initrd && \
    find . | cpio -H newc -o | gzip -c -1 > /output/initrd
