ARG REPO
ARG TAG
FROM ${REPO}/k3os-tar:${TAG} as tar

ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG}
ARG ARCH

COPY --from=tar /output/userspace.tar /output/k3os-rootfs-${ARCH}.tar
RUN gzip /output/k3os-rootfs-${ARCH}.tar
