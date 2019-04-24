### BASE ###
FROM alpine:3.8 as base
RUN apk -U add findutils iptables ipset bash eudev tinyssh e2fsprogs mdadm rsync ca-certificates jq iproute2 nfs-utils parted coreutils libseccomp conntrack-tools

FROM base as rio-build
COPY --from=base /bin /usr/src/image/bin/
COPY --from=base /lib /usr/src/image/lib/
COPY --from=base /sbin /usr/src/image/sbin/
COPY --from=base /etc/ssl/certs/ca-certificates.crt /usr/src/image/etc/ssl/certs/ca-certificates.crt
COPY --from=base /etc/terminfo /usr/src/image/etc/terminfo
COPY --from=base /usr /usr/src/image/usr/

WORKDIR /usr/src/image

RUN rm -rf  usr/bin/iconv \
            usr/bin/scanelf \
            usr/bin/ssl_client \
            usr/bin/pkgconf \
            usr/bin/getent \
            usr/bin/locate \
            usr/bin/updatedb \
            usr/bin/c_rehash \
            usr/bin/getconf \
            usr/etc \
            usr/include \
            usr/lib/bash \
            usr/lib/krb5 \
            usr/lib/pkgconfig \
            usr/lib/tc \
            usr/libexec \
            usr/local \
            usr/sbin/nfsiostat \
            usr/sbin/rpc.gssd \
            usr/sbin/nfsidmap \
            usr/sbin/blkmapd \
            usr/sbin/conntrackd \
            usr/sbin/nfct \
            usr/sbin/nfsstat \
            usr/sbin/mountstats \
            usr/sbin/setcap \
            usr/sbin/exportfs \
            usr/sbin/update-ca-certificates \
            usr/sbin/capsh \
            usr/sbin/getcap \
            usr/sbin/ipset \
            usr/sbin/rpcdebug \
            usr/sbin/start-statd \
            usr/sbin/getpcaps \
            usr/sbin/sm-notify \
            usr/share/aclocal \
            usr/share/apk \
            usr/share/ca-certificates \
            usr/share/man \
            usr/share/misc && \
            find usr/share/terminfo -type f -exec rm {} \; && \
            ln -s xterm-color usr/share/terminfo/x/xterm-256color && \
            rmdir usr/share/terminfo/* || true
RUN rm -rf  bin/sh \
            lib/apk \
            lib/mdev \
            sbin/ss \
            sbin/routel \
            sbin/*-compat* \
            sbin/genl \
            sbin/lnstat \
            sbin/ifstat \
            sbin/mkmntdirs \
            sbin/nfsdcltrack \
            sbin/rtacct \
            sbin/nstat \
            sbin/routef \
            sbin/apk \
            sbin/tc \
            sbin/ifcfg \
            sbin/setup-udev \
            sbin/rtpr \
            sbin/osd_login \
            sbin/bridge \
            sbin/rtmon && \
    ln -s bash bin/sh && \
    mkdir -p lib/modules

RUN mv sbin/* bin/ && \
    rmdir sbin && \
    ln -s bin sbin

RUN mkdir lib2 && \
    mv usr/lib/* lib2/ && \
    mv lib2/* lib/ && \
    mv usr/bin/* bin/ && \
    mv usr/sbin/* bin/ && \
    mv usr/share . && \
    rm -rf usr lib2 && \
    for i in $(ls -l bin | grep usr/bin/coreutils | awk '{print $(NF-2)}'); do \
        rm bin/$i && ln -s coreutils bin/$i \
    ;done && \
    find -L bin -type l -exec rm {} \; -print

RUN apk add upx && \
    upx $(find bin -type f -executable \! -name coreutils) || true

RUN echo '#### LAYOUT #####' && \
    find -type d && \
    echo '#### BIN #####' && \
    find bin -type f -executable && \
    du -x -s -h

RUN tar cvf ../rootfs.tar * && \
    ls -la ../rootfs.tar

CMD ["sh"]

### BUILD IMAGE ###
FROM golang:1.11-alpine AS gobuild
RUN apk -U add git gcc linux-headers musl-dev make libseccomp libseccomp-dev bash
RUN rm -f /bin/sh && ln -s /bin/bash /bin/sh

### ASSEMBLE IMAGE ###
FROM gobuild

RUN apk add -U squashfs-tools

COPY --from=rio-build /usr/src/rootfs.tar /usr/src/rootfs.tar
RUN mkdir /usr/src/image && \
    tar xf /usr/src/rootfs.tar -C /usr/src/image

COPY image/init /usr/src/image/init

RUN mksquashfs /usr/src/image main.squashfs
