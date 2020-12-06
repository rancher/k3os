# package/packer/hetzner

## Setup

Select a project in Hetzner Cloud console and create an API token:

Security --> API Tokens --> Generate API token (Read & Write permisions)

Copy token

## Build

```shell script
export HCLOUD_TOKEN="TOKEN" # Token from Setup
packer build template.json
```
