ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG}

ARG ARCH
ENV ARCH ${ARCH}
ENV VERSION v0.8.1
RUN mkdir -p /output && \
    curl -o /output/install.sh -fL https://raw.githubusercontent.com/rancher/k3s/${VERSION}/install.sh && \
    chmod +x /output/install.sh
RUN INSTALL_K3S_VERSION=${VERSION} INSTALL_K3S_SKIP_START=true INSTALL_K3S_BIN_DIR=/output /output/install.sh
RUN echo "${VERSION}" > /output/version
