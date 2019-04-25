# k3OS
k3OS is a linux distribution designed to remove as much as possible
OS maintaince in a Kubernetes cluster.  It is specifically designed to only
have what is need to run [k3s](https://github.com/rancher/k3s). Additionally
the OS is designed to be managed by kubectl once a cluster is bootstrapped.
Nodes only need to join a cluster and then all aspects of the OS can be managed
from Kubernetes. Both k3OS and k3s upgrades are handled by k3OS.

## Quick Start
Download the ISO from the latest [release](https://github.com/rancher/k3os/releases) and run
in VMware, VirtualBox, or KVM.  The server will automatically start a single node kubernetes cluster. 
Log in with the user `rancher` and run `kubectl`.  This is a "live install" running from the ISO media 
and changes will not persist after reboot. 

To copy k3os to local disk, after logging in as `rancher` run `sudo os-config`. Then remove the ISO 
from the virtual machine and reboot. 

Live install (boot from ISO) requires at least 1GB of RAM. Local install requires 512MB RAM.


# k3OS
k3OS is a Linux distribution designed to remove as much OS maintenance as
possible in a Kubernetes cluster. It is specifically designed to only have what
is needed to run [k3s](https://github.com/rancher/k3s). Additionally the OS is
designed to be managed by kubectl once a cluster is bootstrapped. Nodes only
need to join a cluster and then all aspects of the OS can be managed from
Kubernetes. Both k3OS and k3s upgrades are handled by k3OS.

## Quick Start
Download the ISO from the latest
[release](https://github.com/rancher/k3os/releases) and boot it on VMware,
VirtualBox, or KVM. The server will automatically start a single node cluster.
Log in with the user `rancher` to run kubectl.

## Configuration
All configuration is done through a single cloud-init style config file that is
either packaged in the image, downloaded though cloud-init or managed by
Kubernetes.

More docs to come

## Images
At the moment k3OS will not be shipping official images for the various clouds
providers or other platforms. Instead we will release documentation on how to
automate the installation to master images.

## ARM64 and ARMv7
k3OS will officially ship an ARM64 kernel for server class ARM64 machines. For
smaller ARM64 SBCs and ARMv7, k3OS will support a "Bring Your Own Kernel"
approach. Basically you should be able to use an existing kernel with k3OS and
k3OS wil manage everything except the kernel.

## Management from Kubernetes
Still in development but all configuration of the cluster and nodes will be
accessible from kubectl as Custom Resources. Upgrades of the OS and k3s will be
orchestrated by Kubernetes based on the upgrade policy and desired version.

## License
Copyright (c) 2014-2019 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
