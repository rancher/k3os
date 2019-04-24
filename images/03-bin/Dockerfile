ARG REPO
ARG TAG
FROM ${REPO}/k3os-rootfs:${TAG} as rootfs
RUN echo

ARG REPO
ARG TAG
FROM ${REPO}/k3os-progs:${TAG} as progs
RUN echo

ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG}

ARG REPO
ARG TAG
COPY --from=rootfs /output/rootfs.squashfs /usr/src/ 
COPY --from=progs /output/k3os /output/k3os
RUN echo -n "_sqmagic_" >> /output/k3os && \
    cat /usr/src/rootfs.squashfs >> /output/k3os
