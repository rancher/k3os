ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG} AS base
RUN apk --no-cache add \
    grub-bios \
    open-vm-tools