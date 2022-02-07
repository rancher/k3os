ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG}

ARG ARCH
ENV ARCH ${ARCH}
ENV VERSION v1.23.3+k3s1
ADD https://raw.githubusercontent.com/rancher/k3s/${VERSION}/install.sh /output/install.sh
ENV INSTALL_K3S_VERSION=${VERSION} \
    INSTALL_K3S_SKIP_START=true \
    INSTALL_K3S_BIN_DIR=/output
RUN chmod +x /output/install.sh
RUN /output/install.sh
RUN echo "${VERSION}" > /output/version
