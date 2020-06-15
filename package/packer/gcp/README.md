# package/packer/gcp

## Setup

Configure a Compute Engine Service Account (`account.json`):
- https://www.packer.io/docs/builders/googlecompute/#running-without-a-compute-engine-service-account

Configure `${GCP_PROJECT_ID}`.

## Build AMD64

```shell script
packer build template.json
```
