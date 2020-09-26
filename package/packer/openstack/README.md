# k3OS Packer OpenStack image

These are templates for building a k3OS [OpenStack](https://openstack.org) image for both AMD64 and ARM64 using [Packer](https://www.packer.io).

## Quick start

### Pre-requisites

Configure access and authentication to the target environment as per the official documentation [here](https://docs.openstack.org/python-openstackclient/pike/configuration/index.html).

You'll need an existing network that's capable of routing traffic via the Internet in order to boot your instance.  If your provider has enabled the [appropriate extensions](https://docs.openstack.org/neutron/ussuri/admin/config-auto-allocation.html), this can be quickly performed with the following command:

```
openstack network auto allocated topology create --or-show
```

If your cloud provider doesn't have these features enabled, you'll have to first create your own functioning network topology either manually or using tools such as Terraform.

You'll need the network ID for your network as well as the ID or name of your provider's floating IP pool along with the flavor (size) of image to be used during the build process.  In the example below, 'frankfurter' is equivalent to a VM with 1 vCPU, 2GB RAM and a 20GB disk.

You also need to make sure that the default Security Group has a rule to allow SSH ingress.

### Running

The following example sets the appropriate environment variables and then runs `packer` to validate and build the image.

```
export OS_SOURCE_IMAGE="948f0a55-0b7a-4f8b-ba6e-1461b13e3ea9"
export OS_NETWORK_ID="7dad98ae-65d7-4ec4-b818-73481ff06e3e"
export OS_FLOATING_IP_POOL="internet"
export OS_FLAVOR="frankfurter"
packer validate template.json
packer build template.json
```

> For ARM64, replace `template.json` with `template-arm64.json`.

### Testing

A successful build should complete with the ID of the newly-created image:

```
==> openstack: Deleting temporary keypair: packer_5f16f42b-a7bf-6b2c-3853-180a6c085d88 ...
Build 'openstack' finished.

==> Builds finished. The artifacts of successful builds are:
--> openstack: An image was created: 398801c7-eaa8-46d1-bee2-e53f60b64c3b
```

You should then be able to boot this image, attach a floating IP, and SSH:

```bash
$ openstack server create \
> --image 398801c7-eaa8-46d1-bee2-e53f60b64c3b \
> --network $OS_NETWORK_ID \
> --flavor $OS_FLAVOR \
> --key-name Deadline \
> --wait k3os-test

+-----------------------------+-----------------------------------------------------------+
| Field                       | Value                                                     |
+-----------------------------+-----------------------------------------------------------+
| OS-EXT-STS:power_state      | Running                                                   |
| OS-EXT-STS:vm_state         | active                                                    |
| OS-SRV-USG:launched_at      | 2020-07-21T14:33:20.000000                                |
| addresses                   | auto_allocated_network=10.0.1.53                          |
| config_drive                |                                                           |
| created                     | 2020-07-21T14:32:57Z                                      |
| flavor                      | frankfurter (598a4db4-66a8-49e6-ad8f-d1478afb889b)        |
| id                          | f82eff0c-7adf-496c-88f5-3e40a51fa77a                      |
| image                       | k3OS-v0.10.3-amd64 (398801c7-eaa8-46d1-bee2-e53f60b64c3b) |
| key_name                    | Deadline                                                  |
| name                        | k3os-test                                                 |
| security_groups             | name='default'                                            |
| status                      | ACTIVE                                                    |
| updated                     | 2020-07-21T14:33:20Z                                      |
+-----------------------------+-----------------------------------------------------------+

$ openstack floating ip create $OS_FLOATING_IP_POOL | grep floating_ip_address
| floating_ip_address | 193.16.42.77 |

$ openstack server add floating ip k3os-test 193.16.42.77

$ ssh rancher@193.16.42.77

Warning: Permanently added '193.16.42.77' (ECDSA) to the list of known hosts.
Welcome to k3OS!
```

## Notes

### ARM64 notes

It was developed using ARM64 images on an Ampere Computing eMAG running an OpenStack All-In-One deployed via Kolla and Kolla Ansible.

The image used was a stock Ubuntu 18.04 cloud image for arm64.  The image was loaded into the OpenStack deployment using Terraform from the project located [here](https://github.com/amperecomputing/terraform-openstack-images).


I had to build Packer from source for aarch64 due to the following issue with Packer: [8258](https://github.com/hashicorp/packer/issues/8258)

