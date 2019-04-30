#!/bin/bash
set -e

PROG=$0

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
        umount ${TARGET}/boot || true
        rmdir ${TARGET}/boot || true
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
    echo "Usage: $PROG [--efi] [--msdos] [--config https://.../config.yaml] DEVICE ISO_URL"
    echo ""
    echo "Example: $PROG --efi /dev/vda https://github.com/rancher/k3os/releases/download/v0.2.0-rc3/k3os.iso"
    exit 1
}

if [ "$K3OS_DEBUG" = true ]; then
    set -x
fi

while [ "$#" -gt 0 ]; do
    case $1 in
        --no-format)
            K3OS_INSTALL_NO_FORMAT=true
            ;;
        --msdos)
            K3OS_INSTALL_MSDOS=true
            ;;
        --efi)
            K3OS_INSTALL_EFI=true
            ;;
        --poweroff)
            K3OS_INSTALL_POWER_OFF=true
            ;;
        --config)
            shift 1
            K3OS_INSTALL_CONFIG_URL=$1
            ;;
        -h)
            usage
            ;;
        --help)
            usage
            ;;
        *)
            if [ "$#" != 2 ]; then
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
    echo "Usage: $0 [--efi] [--msdos] [--config https://.../config.yaml] DEVICE ISO_URL"
    echo ""
    echo "Example: $0 --efi /dev/vda https://github.com/rancher/k3os/releases/download/v0.2.0-rc3/k3os.iso"
    exit 1
fi

for i in dd syslinux curl mkfs.ext4 mkfs.vfat fatlabel parted partprobe rsync; do
    if [ ! -x "$(which $i)" ]; then
        MISSING="${MISSING} $i"
    fi
done

if [ -n "${MISSING}" ]; then
    echo "The following required programs are missing for installation: ${MISSING}"
    exit 1
fi

DEVICE=$K3OS_INSTALL_DEVICE
PARTTABLE=${PARTTABLE:-gpt}

if [ "$K3OS_INSTALL_EFI" ]; then
    PARTTABLE=efi
elif [ "$K3OS_INSTALL_MSDOS" ]; then
    PARTTABLE=msdos
fi

if [ ! -b ${DEVICE} ]; then
    echo "You should use an available device. Device ${DEVICE} does not exist."
    exit 1
fi

if [ "${PARTTABLE}" != gpt ] && [ "${PARTTABLE}" != msdos ] && [ "${PARTTABLE}" != "efi" ]; then
    echo "Invalid partition table type ${PARTTABLE}, must be either gpt, msdos, efi"
    exit 1
fi

trap cleanup exit

ISO_DEVICE=$(blkid -L K3OS || true)
if [ -z "${ISO_DEVICE}" ] && [ -n "$K3OS_INSTALL_ISO_URL" ]; then
    TEMP_FILE=$(mktemp k3os.XXXXXXXX.iso)
    get_url ${K3OS_INSTALL_ISO_URL} ${TEMP_FILE}
    ISO_DEVICE=$(losetup --show -f $TEMP_FILE)
    rm -f $TEMP_FILE
fi

if [ -z "${ISO_DEVICE}" ]; then
    echo "#### There is no k3os ISO device"
    exit 1
fi

BOOTFLAG=legacy_boot
MBR_FILE=gptmbr.bin
if [ "${PARTTABLE}" = "msdos" ]; then
    BOOTFLAG=boot
    MBR_FILE=mbr.bin
elif [ "${PARTTABLE}" = "efi" ]; then
    EFI=true
    PARTTABLE=gpt
    BOOTFLAG=esp
    if [ ! -e /sys/firmware/efi ]; then
        echo WARNING: installing EFI on to a system that does not support EFI
        HAS_EFI=false
    fi
fi

if [ "$K3OS_INSTALL_NO_FORMAT" = "true" ]; then
    BOOT=$(blkid -L K3OS_BOOT)
    STATE=$(blkid -L K3OS_STATE)
else
    BOOT_NUM=1
    STATE_NUM=2

    dd if=/dev/zero of=${DEVICE} bs=1M count=1
    parted -s ${DEVICE} mklabel ${PARTTABLE}
    parted -s ${DEVICE} mkpart primary fat32 0% 500MB
    parted -s ${DEVICE} mkpart primary ext4 500MB 1000MB
    parted -s ${DEVICE} set 1 ${BOOTFLAG} on
    partprobe

    PREFIX=${DEVICE}
    if [ ! -e ${PREFIX}${STATE_NUM} ]; then
        PREFIX=${DEVICE}p
    fi

    if [ ! -e ${PREFIX}${STATE_NUM} ]; then
        echo Failed to find ${PREFIX}${STATE_NUM} or ${DEVICE}${STATE_NUM} to format
        exit 1
    fi

    BOOT=${PREFIX}${BOOT_NUM}
    STATE=${PREFIX}${STATE_NUM}

    mkfs.ext4 -F -L K3OS_STATE ${STATE}
    mkfs.vfat -F 32 ${BOOT}
    fatlabel ${BOOT} K3OS_BOOT
fi


TARGET=/run/k3os/target
mkdir -p ${TARGET}
mount ${STATE} ${TARGET}
mkdir -p ${TARGET}/boot
mount ${BOOT} ${TARGET}/boot

DISTRO=/run/k3os/iso
mkdir -p $DISTRO
mount $ISO_DEVICE $DISTRO

rsync -av --exclude boot/isolinux ${DISTRO}/ $TARGET
if [ -n "$STATE_NUM" ]; then
    echo $DEVICE $STATE > $TARGET/k3os/system/growpart
fi

mkdir -p ${TARGET}/boot/EFI
cp -rf ${DISTRO}/boot/isolinux ${TARGET}/boot/syslinux
cp -rf ${DISTRO}/boot/isolinux/efi64 ${TARGET}/boot/EFI/syslinux
cp -rf ${DISTRO}/boot/isolinux/efilinux.cfg ${TARGET}/boot/EFI/syslinux/syslinux.cfg
TTY=$(tty | sed 's!/dev/!!')
if [ "$TTY" != tty1 ] && [ -n "$TTY" ]; then
    sed -i "s!console=tty1!console=tty1 console=${TTY}!g" ${TARGET}/boot/syslinux/syslinux.cfg ${TARGET}/boot/EFI/syslinux/syslinux.cfg
fi

if [ "$EFI" = "true" ]; then
    modprobe efivars 2>/dev/null || true
    if [ "${HAS_EFI}" != "false" ]; then
        efibootmgr -c -d ${DEVICE} -p 1 -l \\EFI\\syslinux\\syslinux.efi -L "SYSLINUX"
    fi
else
    syslinux --directory /syslinux --install ${BOOT}
    dd bs=440 conv=notrunc count=1 if=${TARGET}/boot/syslinux/${MBR_FILE} of=${DEVICE}
fi

if [ -n "$K3OS_INSTALL_CONFIG_URL" ]; then
    get_url "$K3OS_INSTALL_CONFIG_URL" ${TARGET}/k3os/system/config.yaml
    chmod 600 ${TARGET}/k3os/system/config.yaml
fi

if [ -n "$INTERACTIVE" ]; then
    exit 0
fi

if [ "$K3OS_INSTALL_POWER_OFF" = true ] || grep -q 'k3os.mode=install' /proc/cmdline; then
    poweroff -f
else
    echo " * Rebooting system in 5 seconds"
    sleep 5
    reboot -f
fi
