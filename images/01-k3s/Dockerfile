ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG}

ARG ARCH
ENV ARCH ${ARCH}
ENV VERSION v0.5.0
RUN mkdir -p /output && \
    curl -o /output/install.sh -fL https://raw.githubusercontent.com/ibuildthecloud/k3s-dev/5d5352ba1ca742199afa062ca08bf56f40b71d4b/install.sh && \
    chmod +x /output/install.sh
RUN K3S_VERSION=${VERSION} INSTALL_K3S_SKIP_START=true INSTALL_K3S_BIN_DIR=/output /output/install.sh
RUN echo "${VERSION}" > /output/version
