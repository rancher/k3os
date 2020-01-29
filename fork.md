# k3OS
This fork of the upstream k3OS add as default
- MetalLB as default load-balancer
- OpenEBS.io with cStor as storage provider
- Helm (and Tiller) with RBAC to support Helm-Chart package deployments

Please see [README.md] for installation 

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