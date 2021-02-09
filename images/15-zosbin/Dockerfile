### BASE ###
FROM alpine:3.12 as zosbin
ARG ARCH

RUN apk --no-cache add alpine-sdk fuse3-dev hiredis-dev linux-headers

RUN git clone https://github.com/threefoldtech/0-db-fs && \
  cd 0-db-fs && \
  make release

RUN git clone https://github.com/threefoldtech/0-db && \
  cd 0-db/libzdb && \
  make release && \
  cd ../zdbd && \
  make release

