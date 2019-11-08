# k3OS Packer VMware image

This is a working template for building k30s Server and Agent templates using packer in VMware.

It was developed on vCenter 6.7 U3 and ESXi 14320388.

This image does not utilize a cloud-init config and instead utilizes the boot command and networking is done via DHCP.

## Quick start

Download the [latest version](https://github.com/rancher/k3os/releases/latest) version of k3os and copy it to your vCenter datastore.

Assuming that packer and packer-builder-vware-iso are installed you will run the following commands:

packer.io build -var-file k3os-server-variables.json k3os-server.json

Example agent-variable.json

```
{
    "vcenter_server": "vcenter.example.com",
    "vcenter_username": "administrator@vsphere.local",
    "vcenter_password": "VMware123!",
    "vcenter_datastore": "datastore0",
    "vcenter_folder": "Packer_Images",
    "vcenter_host": "esxi.example.com",
    "vcenter_network": "10.0.0.x-24",
    "vcenter_iso_path": "[datastore0] ISOs/k3os-amd64.iso",
    "hostname": "k3os-agent-template",
    "ssh_username": "rancher",
    "rancher_password": "P@$$w0rd1",
    "server_url": "https://10.0.0.50:6443",
    "server_token": "K10cec0c326040384622e1fed081deee3e06::node:2c0e5510ed6d797ea3f1c1"
}
```

## Demo deploying the k3os agent vm from a template

[![asciicast](https://asciinema.org/a/73CNGRExNUaLDGA545Bczg770.svg)](https://asciinema.org/a/73CNGRExNUaLDGA545Bczg770)
