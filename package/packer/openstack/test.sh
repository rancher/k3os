#!/usr/bin/env bash

OS_PACKER_SOURCE_IMAGE_ID='cc568f99-d3b2-489b-b7c6-877e3d6ff1fd'
OS_NETWORKS_ID='11e4411f-e3ff-414e-a216-15b998e03e76'
OS_FLOATING_IP_POOL='public1'
echo $OS_PACKER_SOURCE_IMAGE_ID
echo $OS_NETWORKS_ID
echo $OS_FLOATING_IP_POOL


echo "Validating AMD64 Template"
packer validate template.json
echo "Validating ARM64 Template"
packer validate template-arm64.json

packer build template-arm64.json


