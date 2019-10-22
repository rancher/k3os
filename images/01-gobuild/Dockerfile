FROM golang:1.12-alpine3.9 AS gobuild
RUN apk -U add git gcc linux-headers musl-dev make libseccomp libseccomp-dev bash
COPY gobuild gotest /usr/bin/
RUN rm -f /bin/sh && ln -s /bin/bash /bin/sh
