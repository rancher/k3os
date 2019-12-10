ARG TAG
ARG REPO
FROM ${REPO}/k3os-package:${TAG} AS package

FROM scratch
ENV PATH /k3os/system/k3os/current:/k3os/system/k3s/current
COPY --from=package /output/k3os/system/ /k3os/system/
