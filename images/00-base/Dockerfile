### BASE ###
FROM alpine:3.9 as base
ARG ARCH
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
    grub-efi \
    haveged \
    iproute2 \
    iptables \
    jq \
    kbd-bkeymaps \
    logrotate \
    nfs-utils \
    open-iscsi \
    openrc \
    openssh-client \
    openssh-server \
    parted \
    procps \
    rsync \
    strace \
    sudo \
    tar \
    util-linux \
    vim \
    xz
RUN if [ "$ARCH" == "amd64" ]; then \
        apk add open-vm-tools grub-bios \
    ;fi
RUN cp /etc/apk/repositories /etc/apk/repositories.orig && \
    echo 'http://dl-3.alpinelinux.org/alpine/edge/testing' >> /etc/apk/repositories && \
    apk -U add efibootmgr && \
    mv /etc/apk/repositories.orig /etc/apk/repositories && \
    apk update
