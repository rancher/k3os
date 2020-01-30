# k3OS baremetal
This fork of the upstream k3OS add as default
- [MetalLB layer 2](https://metallb.universe.tf/concepts/layer2/) as default load-balancer
- [OpenEBS.io](https://github.com/openebs/openebs) with [cStor](https://github.com/openebs/cstor) as storage provider
- Helm 3 (without Tiller) with RBAC to support Helm-Chart package deployments
- Preconfigure using response/config-file within the iso

Please see [k3OS README](README.md) for installation 

Readings
- [Deploy k3os and openebs](https://medium.com/@fromprasath/deploy-k3s-cluster-on-k3os-and-use-openebs-as-persistent-storage-provisioner-3db229c0acf8)
- [K3OS with MetalLB and Dashboad](https://mindmelt.nl/mindmelt.nl/2019/04/08/k3s-kubernetes-dashboard-load-balancer/)

# Preconfigure k3OS config.yaml
The ```/k3os/system/config.yaml``` file is reserved for the system installation and should not be modified on a running system.
At runtime, the config can be changed by creating/modifying the ```/var/lib/rancher/k3os/config.yaml``` and ```/var/lib/rancher/k3os/config.d/*```

## Configure k3s (Kubernetes) within k3OS
All Kubernetes configuration is done by configuring k3s. This is primarily done through environment and k3s_args keys in config.yaml.
Default Environment variable can be added by modifying [environment variable](overlay/etc/environment).
The write_files key can be used to populate the /var/lib/rancher/k3s/server/manifests folder with apps you'd like to deploy on boot.
Any file found in [k3s /var/lib/rancher/k3s/server/manifests](overlay/share/rancher/k3s/server/manifests) will automatically be deployed to Kubernetes in a manner similar to kubectl apply.
It is also possible to deploy Helm charts. K3s supports a CRD controller for installing charts. 
See [k3s yaml manifests](https://github.com/rancher/k3s/tree/master/manifests) as examples.

see [k3s advanced options](https://rancher.com/docs/k3s/latest/en/advanced/) for configuring Helm, MetalLB, OpenEBS.

## SSH Public-Key
- Generate a public key fingerprint for you (SSH client on Windows: https://www.ssh.com/ssh/putty/windows/puttygen#running-puttygen)
- Add your public-key to [K3OS configuration file baked into the iso](images/07-iso/config.yaml)
```
  ssh_authorized_keys:
    - ssh-rsa YOURSSHPUBLICKEY rancher@myhostname
  hostname: myhostname
``` 

# Deploy in Windows 10 Hyper-V for dry-run tests
Create a V-2 virtual machine loading your iso file
Take care to:
- Don't use secure-boot.
- Use Default-Network

### Setup SSH access

### Access from Host-OS:
- In K3os-VM, run ```$ ifconfig```, and read under eth0 the ip-address
- In K3os-VM, run ```$ hostname
k3os-1962```
- Generate a public key fingerprint for your client (https://www.ssh.com/ssh/putty/windows/puttygen#running-puttygen)
  Output: 
  
  sudo vi /var/lib/rancher/k3os/config.yaml
  
  Add following ```
  ssh_authorized_keys:
  - ssh-rsa SHA256:d3Ykm9jBz1n/socoIZeFeL7c1PpnlOG7W7aR0CxagC8 rancher@k3os-1962
  ``` 

# Fork Maintainance
## Configure your fork of k3OS
1. Configure your remote fork
```
$ git remote -v
origin  https://github.com/ccuz/k3os.git (fetch)
origin  https://github.com/ccuz/k3os.git (push)
```
1. Add k3OS as upstream project
```
$ git remote add upstream https://github.com/rancher/k3os.git
```
1. Verify
```
$ git remote -v
origin  https://github.com/ccuz/k3os.git (fetch)
origin  https://github.com/ccuz/k3os.git (push)
upstream        https://github.com/rancher/k3os.git (fetch)
upstream        https://github.com/rancher/k3os.git (push)
```

## Merge upstream back into your fork of k3OS
1. Fetch the upstream changes ```$ git fetch upstream```
1. Create a new local branch to merge the upstream changes into your fork ```$ git checkout -b merge-upstream-in-my-fork```
1. Merge the upstream changes into my branch ```$ git merge upstream/master```
1. Check build and make a PR to merge your 'merge-upstream-in-my-fork' into your origin/master