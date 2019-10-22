ARG REPO
ARG TAG
FROM ${REPO}/k3os-gobuild:${TAG} as gobuild

ENV LINUXKIT 63afed7b131eef89d3a36b077d090a8e2b790b7a

FROM gobuild as linuxkit
RUN mkdir -p /output && \
    mkdir -p $GOPATH/src/github.com/linuxkit && \
    cd $GOPATH/src/github.com/linuxkit && \
    git clone https://github.com/ibuildthecloud/linuxkit.git && \
    cd linuxkit && \
    git checkout $LINUXKIT
WORKDIR $GOPATH/src/github.com/linuxkit/linuxkit
RUN cd pkg/init/cmd/rc.init && \
    gobuild -o /output/rc.init
RUN cd pkg/metadata && \
    gobuild -o /output/metadata

FROM gobuild as k3os
RUN apk add rsync
COPY /cmd/ $GOPATH/src/github.com/rancher/k3os/cmd/
COPY /pkg/ $GOPATH/src/github.com/rancher/k3os/pkg/
COPY /main.go $GOPATH/src/github.com/rancher/k3os/
COPY /vendor/ $GOPATH/src/github.com/rancher/k3os/vendor/
WORKDIR $GOPATH/src/github.com/rancher/k3os
RUN gobuild -o /output/k3os && \
    gobuild -o /output/k3os-config ./cmd/config && \
    gobuild -o /output/ccapply ./cmd/ccapply && \
    cd pkg/test && \
    find . -mindepth 1 -type d -exec sh -xc 'gotest -c -o {}/package.test {} || true' \; && \
    rsync -av --exclude '*.go' ../test/ /output/test/ && \
    rsync -av ../testdata/ /output/test/

FROM gobuild
COPY --from=linuxkit /output/ /output/
COPY --from=k3os /output/ /output/
WORKDIR /output
RUN set -x \
 && git clone --branch v0.7.0 https://github.com/ahmetb/kubectx.git \
 && chmod -v +x kubectx/kubectx kubectx/kubens
