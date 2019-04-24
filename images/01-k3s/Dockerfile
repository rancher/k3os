ARG REPO
ARG TAG
FROM ${REPO}/k3os-base:${TAG}

ARG ARCH
ENV VERSION v0.4.0
RUN SUFFIX="-${ARCH}" && \
    if [ "$ARCH" = "arm" ]; then SUFFIX="-armhf"; fi && \
    if [ "$ARCH" = "amd64" ]; then SUFFIX=""; fi && \
    mkdir /output && \
    echo https://github.com/rancher/k3s/releases/download/${VERSION}/k3s${SUFFIX} > /output/url && \
    curl -fL $(cat /output/url) -o /output/k3s && \
    chmod +x /output/k3s
RUN echo "${VERSION}" > /output/version
