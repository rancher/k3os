ARG REPO
ARG TAG
FROM ${REPO}/k3os-rootfs:${TAG} as rootfs

ARG REPO
ARG TAG
FROM ${REPO}/k3os-progs:${TAG} as progs

ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG}

COPY --from=rootfs /output/rootfs.squashfs /usr/src/
COPY install.sh /output/k3os-install.sh
COPY --from=progs /output/k3os /output/k3os
RUN echo -n "_sqmagic_" >> /output/k3os
RUN cat /usr/src/rootfs.squashfs >> /output/k3os
