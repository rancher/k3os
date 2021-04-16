ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG} AS base
RUN apk --no-cache add \
    grub-bios \
    open-vm-tools \
    open-vm-tools-deploypkg \
    open-vm-tools-guestinfo \
    open-vm-tools-static \
    open-vm-tools-vmbackup
