#!/usr/bin/env sh

: ${K3OS_UPDATE_CHANNEL_URL:=${K3OS_UPGRADE_CHANNEL="github-releases://rancher/k3os"}}
: ${K3OS_UPDATE_CHANNEL_NAME:=${K3OS_UPDATE_CHANNEL_URL%%:*}}
: ${K3OS_UPDATE_CHANNEL_NAMESPACE:="k3os-system"}
: ${K3OS_UPDATE_CHANNEL_VERSION:=${K3OS_VERSION}}
: ${K3OS_UPDATE_CHANNEL_CONCURRENCY:=1}

if [ -z "${K3OS_UPDATE_CHANNEL_VERSION}" ] && [ -e /etc/os-release ]; then
    . /etc/os-release
    if [ "$ID" = "k3os" ]; then
        K3OS_UPDATE_CHANNEL_VERSION=$VERSION_ID
    fi
fi

cat << EOF
---
apiVersion: k3os.cattle.io/v1
kind: UpdateChannel
metadata:
  name: ${K3OS_UPDATE_CHANNEL_NAME}
  namespace: ${K3OS_UPDATE_CHANNEL_NAMESPACE}
spec:
  concurrency: ${K3OS_UPDATE_CHANNEL_CONCURRENCY}
  version: ${K3OS_UPDATE_CHANNEL_VERSION:="latest"}
  url: ${K3OS_UPDATE_CHANNEL_URL}
EOF
