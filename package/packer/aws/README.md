# package/packer/aws

## Setup

First, setup your AWS environment for Packer as per your preferred method:
- https://www.packer.io/docs/builders/amazon/#static-credentials
- https://www.packer.io/docs/builders/amazon/#environment-variables
- https://www.packer.io/docs/builders/amazon/#shared-credentials-file

## Build AMD64

```shell script
packer build template.json
```

## Build ARM64

```shell script
packer build template-arm64.json
```
