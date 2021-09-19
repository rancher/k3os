FROM golang:1.16-alpine3.14 AS gobuild
RUN apk -U add git gcc linux-headers musl-dev make libseccomp libseccomp-dev bash
COPY gobuild /usr/bin/
RUN rm -f /bin/sh && ln -s /bin/bash /bin/sh
WORKDIR /output