# k3OS Packer OpenStack image

This is a template for building a K3OS OpenStack Image using packer.  

It was developed using ARM64 images on an Ampere Computing eMAG running an OpenStack All-In-One deployed via Kolla and Kolla Ansible.  

The image used was a stock Ubuntu 18.04 cloud image for arm64.  The image was loaded into the openstack deployment using Terrfaform from the project located [here](https://github.com/amperecomputing/terraform-openstack-images).


## Quick start

Assuming the OpenStack deployment was deployed with Kolla and Kolla-ansible. And you obtained the id for the Ubuntu Image 

```
source /etc/kolla/admin-openrc.sh
export OS_SOURCE_IMAGE="abcdefgh-1234-5678-abcd-efghijklmnop"
export OS_NETWORKS_ID="12345678-abcd-efgh-ijkl-123456789abc"
export OS_FLOATING_IP_POOL="public1"
export OS_FLAVOR="2"
packer validate template-arm64.json
packer build template-arm64.json
```

## Notes

* I had to build packer from source for aarch64 due to the following issue with packer: [8258](https://github.com/hashicorp/packer/issues/8258) 


## References

* [https://github.com/hashicorp/packer/issues/8258](https://github.com/hashicorp/packer/issues/8258)
