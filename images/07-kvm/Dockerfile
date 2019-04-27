ARG REPO
ARG TAG
FROM ${REPO}/k3os-iso:${TAG} as iso

ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG} as base
ARG TAG
RUN apk add qemu-system-x86_64 libvirt qemu-img ovmf
RUN ln -s /usr/bin/qemu-system-* /usr/bin/qemu-system
RUN qemu-img create -f qcow2 /hd.img 40G

COPY --from=iso /output/ /output/

COPY run-kvm.sh /usr/bin/
CMD ["run-kvm.sh"]
