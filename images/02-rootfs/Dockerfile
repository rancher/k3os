ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG} as base

ARG REPO
ARG TAG
FROM ${REPO}/k3os-progs:${TAG} as progs

ARG REPO
ARG TAG
FROM ${REPO}/k3os-k3s:${TAG} as k3s

FROM base as k3os-build
ARG VERSION
ARG ARCH
RUN apk add squashfs-tools
COPY --from=base /bin /usr/src/image/bin/
COPY --from=base /lib /usr/src/image/lib/
COPY --from=base /sbin /usr/src/image/sbin/
COPY --from=base /etc /usr/src/image/etc/
COPY --from=base /usr /usr/src/image/usr/

# Fix up more stuff to move everything to /usr
RUN cd /usr/src/image && \
    for i in usr/*; do \
        if [ -e $(basename $i) ]; then \
            tar cvf - $(basename $i) | tar xvf - -C usr && \
            rm -rf $(basename $i) \
        ;fi && \
        mv $i . \
   ;done && \
   rmdir usr

# Fix coreutils links
RUN cd /usr/src/image/bin && \
    mv coreutils coreutils.save && \
    find -xtype l -exec rm {} \; -exec ln -s coreutils {} \; && \
    mv coreutils.save coreutils

# Fix sudo
RUN chmod +s /usr/src/image/bin/sudo

# Add empty dirs to bind mount
RUN mkdir -p /usr/src/image/lib/modules

# setup /usr/local
RUN rm -rf /usr/src/image/local && \
    ln -s /var/local /usr/src/image/local

# cleanup files hostname/hosts
RUN rm -rf \
    /usr/src/image/etc/hosts \
    /usr/src/image/etc/hostname \
    /usr/src/image/etc/alpine-release \
    /usr/src/image/etc/apk \
    /usr/src/image/etc/ca-certificates* \
    /usr/src/image/etc/os-release \
    /usr/src/image/etc/ssl/* \
    && ln -s /usr/lib/os-release /usr/src/image/etc/os-release && \
    mkdir -p /usr/src/image/etc/ssl/certs/ && \
    cp -rf /etc/ssl/certs/ca-certificates.crt /usr/src/image/etc/ssl/certs

RUN rm -rf \
    /usr/src/image/usr/sbin/apk \
    /usr/src/image/usr/include \
    /usr/src/image/usr/lib/apk \
    /usr/src/image/usr/lib/pkgconfig \
    /usr/src/image/usr/lib/systemd \
    /usr/src/image/usr/lib/udev \
    /usr/src/image/usr/share/apk \
    /usr/src/image/usr/share/applications \
    /usr/src/image/usr/share/ca-certificates \
    /usr/src/image/usr/share/icons \
    /usr/src/image/usr/share/mkinitfs \
    /usr/src/image/usr/share/vim/vim81/spell \
    /usr/src/image/usr/share/vim/vim81/tutor \
    /usr/src/image/usr/share/vim/vim81/doc

COPY --from=k3s /output/install.sh /usr/src/image/libexec/k3os/k3s-install.sh
COPY --from=progs /output/rc.init /usr/src/image/bin/rc.init
COPY --from=progs /output/k3os-config /usr/src/image/libexec/k3os/k3os-config
COPY --from=progs /output/ccapply /usr/src/image/sbin/ccapply
COPY --from=progs /output/metadata /usr/src/image/sbin/metadata

COPY overlay/ /usr/src/image/
COPY install.sh /usr/src/image/libexec/k3os/install
RUN sed -i -e "s/%VERSION%/${VERSION}/g" -e "s/%ARCH%/${ARCH}/g" /usr/src/image/lib/os-release
RUN mkdir -p /output && \
    mksquashfs /usr/src/image /output/rootfs.squashfs
