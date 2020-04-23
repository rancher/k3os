ARG REPO
ARG TAG
FROM ${REPO}/k3os-iso:${TAG}

RUN apk add qemu-system-$(uname -m) ovmf

RUN case $(uname -m) in \
        x86_64) \
            EFI_HYBRID="k3os.install.efi_hybrid=true" && \
            DEVICE="/dev/sda" && \
            CONSOLE="ttyS0" \
            ;; \
        aarch64) \
            EFI_HYBRID="" && \
            DEVICE="/dev/sda" && \
            CONSOLE="ttyAMA0" \
            ;; \
    esac && \
    sed -i -e "s#k3os\\.mode=live#k3os.mode=install k3os.install.silent=true ${EFI_HYBRID} k3os.install.power_off=true k3os.install.device=${DEVICE}#g;s/^set timeout=10\$/set timeout=1/g;s/console=ttyS0 console=tty1/console=tty1 console=${CONSOLE}/g" /usr/src/iso/boot/grub/grub.cfg && \
    mkdir -p /usr/src/qemu && \
    grub-mkrescue -o /usr/src/qemu/k3os-autoinstall.iso /usr/src/iso/. -- -volid K3OS -joliet on

RUN set -o pipefail && \
    case $(uname -m) in \
        x86_64) \
            QEMU_OPTIONS="-drive file=/usr/share/OVMF/OVMF_CODE.fd,if=pflash,format=raw,unit=0,readonly=on" && \
            cp /usr/share/OVMF/OVMF_VARS.fd /output/k3os-qemu-uefi_VARS.raw \
            ;; \
        aarch64) \
            QEMU_OPTIONS="-drive file=/usr/src/qemu/QEMU_EFI.raw,if=pflash,format=raw,unit=0,readonly=on -machine virt -cpu cortex-a57" && \
            cp /usr/share/OVMF/QEMU_VARS.fd /output/k3os-qemu-uefi_VARS.raw && \
            cp /usr/share/OVMF/QEMU_EFI.fd /usr/src/qemu/QEMU_EFI.raw && \
            dd if=/dev/zero of=/usr/src/qemu/QEMU_EFI.raw bs=1c count=1 seek=67108863 && \
            cp /usr/share/OVMF/QEMU_VARS.fd /output/k3os-qemu-uefi_VARS.raw && \
            dd if=/dev/zero of=/output/k3os-qemu-uefi_VARS.raw bs=1c count=1 seek=67108863 \
            ;; \
    esac && \
    qemu-img create -f qcow2 /output/k3os-qemu.qcow2 750M && \
    timeout 7200 qemu-system-$(uname -m) \
        -drive file=/output/k3os-qemu-uefi_VARS.raw,if=pflash,format=raw,unit=1 \
        ${QEMU_OPTIONS} \
        -hda /output/k3os-qemu.qcow2 \
        -cdrom /usr/src/qemu/k3os-autoinstall.iso \
        -serial stdio \
        -m 2048M \
        2>&1 | tee /usr/src/qemu/k3os-qemu.log