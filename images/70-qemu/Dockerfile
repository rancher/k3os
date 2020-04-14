ARG REPO
ARG TAG
FROM ${REPO}/k3os-package:${TAG} as package

ARG REPO
ARG TAG
FROM debian:stable-slim

RUN apt-get update && apt-get --yes --no-install-recommends install libguestfs-tools linux-image-amd64

RUN curl -L http://dl-cdn.alpinelinux.org/alpine/v3.11/releases/x86_64/alpine-minirootfs-3.11.5-x86_64.tar.gz -o /alpine-minirootfs-3.11.5-x86_64.tar.gz

COPY --from=package /output/ /usr/src/qemu/
COPY config.yaml /usr/src/qemu/k3os/system/

COPY create-qemu-disk.sh /

RUN sh /create-qemu-disk.sh