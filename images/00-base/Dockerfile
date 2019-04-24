### BASE ###
FROM alpine:3.9 as base
RUN apk -U add \
    bash \
    bash-completion \
    blkid \
    busybox-initscripts \
    ca-certificates \
    connman \
    conntrack-tools \
    coreutils \
    curl \
    dbus \
    dosfstools \
    e2fsprogs \
    e2fsprogs-extra \
    findutils \
    haveged \
    iproute2 \
    iptables \
    jq \
    libseccomp \
    logrotate \
    nfs-utils \
    open-iscsi \
    open-vm-tools \
    openrc \
    openssh-client \
    openssh-server \
    parted \
    procps \
    rsync \
    strace \
    squashfs-tools \
    sudo \
    syslinux \
    tar \
    util-linux \
    vim \
    wget \
    xz
RUN cp /etc/apk/repositories /etc/apk/repositories.orig && \
    echo 'http://dl-3.alpinelinux.org/alpine/edge/testing' >> /etc/apk/repositories && \
    apk -U add efibootmgr && \
    mv /etc/apk/repositories.orig /etc/apk/repositories && \
    apk update
