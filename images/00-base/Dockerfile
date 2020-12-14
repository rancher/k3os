### BASE ###
FROM alpine:3.12 as base
ARG ARCH
RUN apk --no-cache add \
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
    dmidecode \
    dosfstools \
    e2fsprogs \
    e2fsprogs-extra \
    efibootmgr \
    eudev \
    findutils \
    grub-efi \
    haveged \
    htop \
    hvtools \
    iproute2 \
    iptables \
    irqbalance \
    iscsi-scst \
    jq \
    kbd-bkeymaps \
    lm-sensors \
    logrotate \
    lsscsi \
    lvm2 \
    lvm2-extra \
    mdadm \
    mdadm-misc \
    mdadm-udev \
    multipath-tools \
    ncurses \
    ncurses-terminfo \
    nfs-utils \
    open-iscsi \
    openrc \
    openssh-client \
    openssh-server \
    parted \
    procps \
    qemu-guest-agent \
    rng-tools \
    rsync \
    strace \
    smartmontools \
    sudo \
    tar \
    tzdata \
    util-linux \
    vim \
    wireguard-tools \
    wpa_supplicant \
    xfsprogs \
    xz \
 && mv -vf /etc/conf.d/qemu-guest-agent /etc/conf.d/qemu-guest-agent.orig \
 && mv -vf /etc/conf.d/rngd             /etc/conf.d/rngd.orig \
 && mv -vf /etc/conf.d/udev-settle      /etc/conf.d/udev-settle.orig \
# replicate the default "no idea, friend" behavior of virt-what
 && touch /usr/sbin/virt-what \
 && chmod +x /usr/sbin/virt-what
