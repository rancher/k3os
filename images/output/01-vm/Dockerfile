ARG REPO
ARG TAG
FROM ${REPO}/k3os-iso:${TAG} as iso

ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG}
ARG ARCH

COPY --from=iso /output/k3os-vm.iso /output/k3os-vm-${ARCH}.iso
