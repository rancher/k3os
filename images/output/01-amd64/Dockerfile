ARG REPO
ARG TAG
FROM ${REPO}/k3os-tar:${TAG} as tar

ARG REPO
ARG TAG
FROM ${REPO}/k3os-iso:${TAG} as iso

ARG REPO
ARG TAG
FROM ${REPO}/k3os-kernel:${TAG} as kernel

ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG}

COPY --from=kernel /output/vmlinuz /output/k3os-vmlinuz-amd64
COPY --from=kernel /output/initrd /output/k3os-initrd-amd64
COPY --from=kernel /output/kernel.squashfs /output/k3os-kernel-amd64.squashfs
COPY --from=kernel /output/version /output/k3os-kernel-version-amd64
COPY --from=iso /output/k3os.iso /output/k3os-amd64.iso

COPY --from=tar /output/ /usr/src/tar
COPY --from=kernel /output/initrd /usr/src/initrd/initrd

RUN mv /usr/src/tar/userspace.tar /output/k3os-rootfs-amd64.tar && \
    gzip /output/k3os-rootfs-amd64.tar
