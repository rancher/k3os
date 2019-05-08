ARG REPO
ARG TAG
FROM ${REPO}/k3os-tar:${TAG} as tar

ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG}
ARG ARCH

COPY --from=tar /output/ /usr/src/tar

RUN mkdir -p /output && \
    gzip -c /usr/src/tar/userspace.tar > /output/k3os-rootfs-${ARCH}.tar.gz
