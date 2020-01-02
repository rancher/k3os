#!/bin/bash
set -e

PROG=$0
PROGS="dd curl mkfs.ext4 mkfs.vfat fatlabel parted partprobe grub-install"
DISTRO=/run/k3os/iso

if [ "$K3OS_DEBUG" = true ]; then
    set -x
fi

get_url()
{
    FROM=$1
    TO=$2
    case $FROM in
        http*)
            curl -o $TO -fL ${FROM}
            ;;
        *)
            cp -f $FROM $TO
            ;;
    esac
}

cleanup2()
{
    if [ -n "${TARGET}" ]; then
        umount ${TARGET}/boot/efi || true
        umount ${TARGET} || true
    fi

    losetup -d ${ISO_DEVICE} || true
    umount $DISTRO || true
}

cleanup()
{
    EXIT=$?
    cleanup2 2>/dev/null || true
    return $EXIT
}

usage()
{
    echo "Usage: $PROG [--force-efi] [--debug] [--tty TTY] [--poweroff] [--takeover] [--no-format] [--config https://.../config.yaml] DEVICE ISO_URL"
    echo ""
    echo "Example: $PROG /dev/vda https://github.com/rancher/k3os/releases/download/v0.8.0/k3os.iso"
    echo ""
    echo "DEVICE must be the disk that will be partitioned (/dev/vda). If you are using --no-format it should be the device of the K3OS_STATE partition (/dev/vda2)"
    echo ""
    echo "The parameters names refer to the same names used in the cmdline, refer to README.md for"
    echo "more info."
    echo ""
    exit 1
}

do_format()
{
    if [ "$K3OS_INSTALL_NO_FORMAT" = "true" ]; then
        STATE=$(blkid -L K3OS_STATE || true)
        if [ -z "$STATE" ] && [ -n "$DEVICE" ]; then
            tune2fs -L K3OS_STATE $DEVICE
            STATE=$(blkid -L K3OS_STATE)
        fi

        return 0
    fi

    dd if=/dev/zero of=${DEVICE} bs=1M count=1
    parted -s ${DEVICE} mklabel ${PARTTABLE}
    if [ "$PARTTABLE" = "gpt" ]; then
        BOOT_NUM=1
        STATE_NUM=2
        parted -s ${DEVICE} mkpart primary fat32 0% 50MB
        parted -s ${DEVICE} mkpart primary ext4 50MB 750MB
    else
        BOOT_NUM=
        STATE_NUM=1
        parted -s ${DEVICE} mkpart primary ext4 0% 700MB
    fi
    parted -s ${DEVICE} set 1 ${BOOTFLAG} on
    partprobe 2>/dev/null || true

    PREFIX=${DEVICE}
    if [ ! -e ${PREFIX}${STATE_NUM} ]; then
        PREFIX=${DEVICE}p
    fi

    if [ ! -e ${PREFIX}${STATE_NUM} ]; then
        echo Failed to find ${PREFIX}${STATE_NUM} or ${DEVICE}${STATE_NUM} to format
        exit 1
    fi

    if [ -n "${BOOT_NUM}" ]; then
        BOOT=${PREFIX}${BOOT_NUM}
    fi
    STATE=${PREFIX}${STATE_NUM}

    mkfs.ext4 -F -L K3OS_STATE ${STATE}
    if [ -n "${BOOT}" ]; then
        mkfs.vfat -F 32 ${BOOT}
        fatlabel ${BOOT} K3OS_GRUB
    fi
}

do_mount()
{
    TARGET=/run/k3os/target
    mkdir -p ${TARGET}
    mount ${STATE} ${TARGET}
    mkdir -p ${TARGET}/boot
    if [ -n "${BOOT}" ]; then
        mkdir -p ${TARGET}/boot/efi
        mount ${BOOT} ${TARGET}/boot/efi
    fi

    mkdir -p $DISTRO
    mount -t iso9660 -o ro $ISO_DEVICE $DISTRO
}

do_copy()
{
    tar cf - -C ${DISTRO} k3os | tar xvf - -C ${TARGET}
    if [ -n "$STATE_NUM" ]; then
        echo $DEVICE $STATE_NUM > $TARGET/k3os/system/growpart
    fi

    if [ -n "$K3OS_INSTALL_CONFIG_URL" ]; then
        get_url "$K3OS_INSTALL_CONFIG_URL" ${TARGET}/k3os/system/config.yaml
        chmod 600 ${TARGET}/k3os/system/config.yaml
    fi

    if [ "$K3OS_INSTALL_TAKE_OVER" = "true" ]; then
        touch ${TARGET}/k3os/system/takeover
    fi
}

install_grub()
{
    if [ "$K3OS_INSTALL_DEBUG" ]; then
        GRUB_DEBUG="k3os.debug"
    fi

    mkdir -p ${TARGET}/boot/grub
    cat > ${TARGET}/boot/grub/grub.cfg << EOF
set default=0
set timeout=10

set gfxmode=auto
set gfxpayload=keep
insmod all_video
insmod gfxterm

menuentry "k3OS Current" {
  search.fs_label K3OS_STATE root
  set sqfile=/k3os/system/kernel/current/kernel.squashfs
  loopback loop0 /\$sqfile
  set root=(\$root)
  linux (loop0)/vmlinuz printk.devkmsg=on console=tty1 $GRUB_DEBUG
  initrd /k3os/system/kernel/current/initrd
}

menuentry "k3OS Previous" {
  search.fs_label K3OS_STATE root
  set root=(\$root)
  linux /k3os/system/kernel/previous/vmlinuz printk.devkmsg=on console=tty1 $GRUB_DEBUG
  initrd /k3os/system/kernel/previous/initrd
}

menuentry "k3OS Rescue Shell" {
  search.fs_label K3OS_STATE root
  set root=(\$root)
  linux /k3os/system/kernel/current/vmlinuz printk.devkmsg=on rescue console=tty1
  initrd /k3os/system/kernel/current/initrd
}
EOF
    if [ -z "${K3OS_INSTALL_TTY}" ]; then
        TTY=$(tty | sed 's!/dev/!!')
    else
        TTY=$K3OS_INSTALL_TTY
    fi
    if [ -e "/dev/$TTY" ] && [ "$TTY" != tty1 ] && [ -n "$TTY" ]; then
        sed -i "s!console=tty1!console=tty1 console=${TTY}!g" ${TARGET}/boot/grub/grub.cfg
    fi

    if [ "$K3OS_INSTALL_NO_FORMAT" = "true" ]; then
        return 0
    fi

    if [ "$K3OS_INSTALL_FORCE_EFI" = "true" ]; then
        GRUB_TARGET="--target=x86_64-efi"
    fi

    grub-install ${GRUB_TARGET} --boot-directory=${TARGET}/boot ${DEVICE}
}

get_iso()
{
    ISO_DEVICE=$(blkid -L K3OS || true)
    if [ -z "${ISO_DEVICE}" ]; then
        for i in $(lsblk -o NAME,TYPE -n | grep -w disk | awk '{print $1}'); do
            mkdir -p ${DISTRO}
            if mount -t iso9660 -o ro /dev/$i ${DISTRO}; then
                ISO_DEVICE="/dev/$i"
                umount ${DISTRO}
                break
            fi
        done
    fi

    if [ -z "${ISO_DEVICE}" ] && [ -n "$K3OS_INSTALL_ISO_URL" ]; then
        TEMP_FILE=$(mktemp k3os.XXXXXXXX.iso)
        get_url ${K3OS_INSTALL_ISO_URL} ${TEMP_FILE}
        ISO_DEVICE=$(losetup --show -f $TEMP_FILE)
        rm -f $TEMP_FILE
    fi

    if [ -z "${ISO_DEVICE}" ]; then
        echo "#### There is no k3os ISO device"
        return 1
    fi
}

setup_style()
{
    if [ "$K3OS_INSTALL_FORCE_EFI" = "true" ] || [ -e /sys/firmware/efi ]; then
        PARTTABLE=gpt
        BOOTFLAG=esp
        if [ ! -e /sys/firmware/efi ]; then
            echo WARNING: installing EFI on to a system that does not support EFI
        fi
    else
        PARTTABLE=msdos
        BOOTFLAG=boot
    fi
}

validate_progs()
{
    for i in $PROGS; do
        if [ ! -x "$(which $i)" ]; then
            MISSING="${MISSING} $i"
        fi
    done

    if [ -n "${MISSING}" ]; then
        echo "The following required programs are missing for installation: ${MISSING}"
        exit 1
    fi
}

validate_device()
{
    DEVICE=$K3OS_INSTALL_DEVICE
    if [ ! -b ${DEVICE} ]; then
        echo "You should use an available device. Device ${DEVICE} does not exist."
        exit 1
    fi
}

create_opt()
{
    mkdir -p "${TARGET}/k3os/data/opt"
}

while [ "$#" -gt 0 ]; do
    case $1 in
        --no-format)
            K3OS_INSTALL_NO_FORMAT=true
            ;;
        --force-efi)
            K3OS_INSTALL_FORCE_EFI=true
            ;;
        --poweroff)
            K3OS_INSTALL_POWER_OFF=true
            ;;
        --takeover)
            K3OS_INSTALL_TAKE_OVER=true
            ;;
        --debug)
            set -x
            K3OS_INSTALL_DEBUG=true
            ;;
        --config)
            shift 1
            K3OS_INSTALL_CONFIG_URL=$1
            ;;
        --tty)
            shift 1
            K3OS_INSTALL_TTY=$1
            ;;
        -h)
            usage
            ;;
        --help)
            usage
            ;;
        *)
            if [ "$#" -gt 2 ]; then
                usage
            fi
            INTERACTIVE=true
            K3OS_INSTALL_DEVICE=$1
            K3OS_INSTALL_ISO_URL=$2
            break
            ;;
    esac
    shift 1
done

if [ -e /etc/environment ]; then
    source /etc/environment
fi

if [ -e /etc/os-release ]; then
    source /etc/os-release

    if [ -z "$K3OS_INSTALL_ISO_URL" ]; then
        K3OS_INSTALL_ISO_URL=${ISO_URL}
    fi
fi

if [ -z "$K3OS_INSTALL_DEVICE" ]; then
    usage
fi

validate_progs
validate_device

trap cleanup exit

get_iso
setup_style
do_format
do_mount
do_copy
install_grub
create_opt

if [ -n "$INTERACTIVE" ]; then
    exit 0
fi

if [ "$K3OS_INSTALL_POWER_OFF" = true ] || grep -q 'k3os.mode=install' /proc/cmdline; then
    poweroff -f
else
    echo " * Rebooting system in 5 seconds (CTRL+C to cancel)"
    sleep 5
    reboot -f
fi
