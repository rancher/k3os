ARG REPO
ARG TAG
FROM ${REPO}/k3os-tar:${TAG} as tar

ARG REPO
ARG TAG
FROM ${REPO}/k3os-kernel:${TAG} as kernel

ARG REPO
ARG TAG
FROM ${REPO}/k3os-progs:${TAG} as tests

ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG} as base
ARG VERSION
ARG ARCH

RUN set -x \
 && apk add \
    grub \
    grub-efi \
    mtools \
    xorriso

COPY --from=kernel /output/ /usr/src/kernel/
COPY --from=tar /output/userspace.tar /usr/src/tars/

WORKDIR /usr/src/kernel
RUN mkdir -p /usr/src/iso/k3os/system/kernel/$(cat version)
RUN cp -vf initrd kernel.squashfs /usr/src/iso/k3os/system/kernel/$(cat version)
RUN ln -vs $(cat version) /usr/src/iso/k3os/system/kernel/current
RUN tar xvf /usr/src/tars/userspace.tar --strip-components=1 -C /usr/src/iso

RUN mkdir -p /output
COPY wrapper /usr/bin/
WORKDIR /usr/src/iso

COPY grub.cfg /usr/src/iso/boot/grub/
COPY config.yaml /usr/src/iso/k3os/system/

RUN grub-mkrescue --xorriso=/usr/bin/wrapper -o /output/k3os.iso . -l -allow-lowercase -V K3OS \
 && test -e /output/k3os.iso # grub-mkrescue does not exit non-zero on failure
COPY --from=tests /output/test/ /usr/src/test/
RUN mkisofs -l -allow-lowercase -iso-level 4 -R -V K3OS_TEST_RUN -o /output/k3os-tests.iso /usr/src/test

CMD ["run-kvm.sh"]
