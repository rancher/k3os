ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG}

ARG ARCH
ENV VERSION v0.4.0
RUN mkdir -p /output && \
    curl -o /output/install.sh -fL https://raw.githubusercontent.com/erikwilson/rancher-k3s/fa673468cefbf38f3b46ddbae58ebfc7cc60b806/install.sh && \
    chmod +x /output/install.sh
RUN K3S_VERSION=${VERSION} INSTALL_K3S_SKIP_START=true INSTALL_K3S_BIN_DIR=/output /output/install.sh
RUN echo "${VERSION}" > /output/version
