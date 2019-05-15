# Build image for GCE

This is an *experimental* image builder for k3os that produces images which can be run on [Google Compute Engine.](https://cloud.google.com/compute/).

You should have [gcloud](https://cloud.google.com/sdk/install) installed and set up.

Set your project:

`gcloud config set project [your project]`

Then to run the image building script:

`bash build-image.sh`

The image building script does the following:

* Creates an Ubuntu VM to run the "take-over" installer path
* Displays installer progress via the VM serial console
* Cleans up image-building VM
* Creates a GCE image after k3os is installed on the disk
* Launches a k3os instance with the new image

After the k3os instance is up and running, use the public IP of the image:

`ssh rancher@[ip-of-vm]`

The image currently is built with the password for the rancher user set to `rancher`.

If you have metadata based ssh-keys for the project or VM, they should be available.

If you want to add a github key as supported by k3os you can use instance metadata.

Edit the VM, and add metadata key of `user-data` with a value conforming to k3os config, for example:

```
ssh_authorized_keys:
- github:ptone
```

You may need to reboot the VM for this instance metadata to be copied to the rancher homedir

## To start a cluster:

Set user-data metadata on the master (assume we call the vm `k3os-master` to:

```
k3os:
  token: myclustersecret
```

And on additional VMs to:

```
k3os:
  server_url: https://k3os-master:6443
  token: myclustersecret
``

After nodes have started and registered, you should be able to see them listed on the master:

```
k3os-master [~]$ kubectl get node
NAME                                  STATUS   ROLES    AGE     VERSION
k3os-master.c.[project-id].internal   Ready    <none>   17m     v1.14.1-k3s.4
node2.c.[project-id].internal         Ready    <none>   15m     v1.14.1-k3s.4
node3.c.[project-id].internal         Ready    <none>   2m27s   v1.14.1-k3s.4
```


## cleanup

```
gcloud compute instances delete k3os --zone=${ZONE} --quiet &
gcloud compute images delete k3os-0-2-1 --quiet 
```