ARG REPO
ARG TAG
FROM ${REPO}/k3os-k3s:${TAG} as k3s

ARG REPO
ARG TAG
FROM ${REPO}/k3os-bin:${TAG} as bin

ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG} as base
ARG VERSION

COPY --from=k3s /output/ /usr/src/k3s
RUN cd /usr/src/k3s && \
    mkdir -p /usr/src/tar/k3os/system/k3s/$(cat version) && \
    cp k3s /usr/src/tar/k3os/system/k3s/$(cat version) && \
    ln -sf $(cat version) /usr/src/tar/k3os/system/k3s/current && \
    mkdir -p /usr/src/tar/sbin && \ 
    ln -sf /k3os/system/k3s/current/k3s /usr/src/tar/sbin/k3s

COPY --from=bin /output/k3os /usr/src/tar/k3os/system/k3os/${VERSION}/k3os
RUN ln -sf ${VERSION} /usr/src/tar/k3os/system/k3os/current && \
    ln -sf /k3os/system/k3os/current/k3os /usr/src/tar/sbin/k3os && \
    ln -sf k3os /usr/src/tar/sbin/init

RUN mkdir -p /output && \
    mv /usr/src/tar /usr/src/${VERSION} && \
    cd /usr/src/ && \
    tar cvf /output/userspace.tar ${VERSION}
