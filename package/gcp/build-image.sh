#!/bin/bash

instance="k3os-image-maker"
zone="${ZONE:-us-central1-a}"
project=$(gcloud config list project --format "value(core.project)" )
image_name="k3os-0-2-1"

echo "Starting Image Builder"
gcloud compute instances create $instance \
--zone=$zone --machine-type=n1-standard-2 \
--image=ubuntu-minimal-1804-bionic-v20190429 \
--image-project=ubuntu-os-cloud \
--no-boot-disk-auto-delete \
--async \
--metadata startup-script='#! /bin/bash
cat >> /tmp/config.yaml <<EOL
k3os:
  datasource: gcp
  password: rancher

EOL

if [ ! -f /var/.k3osstartupinstall ]; then
  sudo apt-get update -y && sudo apt-get install -y dosfstools parted
  wget -P /tmp/ https://raw.githubusercontent.com/rancher/k3os/master/install.sh 
  sudo bash -x /tmp/install.sh --takeover --debug --tty ttyS0 --config /tmp/config.yaml --no-format $(findmnt / -o SOURCE -n) https://github.com/rancher/k3os/releases/download/v0.2.1/k3os-amd64.iso
  touch /var/.k3osstartupinstall
  sync
  sudo reboot
fi'

status=$( gcloud compute instances describe $instance --format='value("status")' --zone=$zone )
if [ $? != 0 ]
then
  echo "instance error"
  exit 1
fi

status_re="PROVISIONING|STAGING|RUNNING"
if ! [[ $status =~ $status_re ]]
then
  echo "Instance already stopping or terminated"
  exit 1
fi

echo -n "Waiting for VM serial port"
until gcloud compute instances describe $instance --format='value("status")' --zone=us-central1-a | grep -q "RUNNING";
do 
  echo -n '.'
  sleep 3; 
done

echo "Observe Build"

serial_cursor=0
cursor_regex='--start=([0-9]+) '
while [ $( gcloud compute instances describe $instance --format='value("status")' --zone=$zone ) == "RUNNING" ]

do
  { cmd_err=$(gcloud compute instances get-serial-port-output $instance --start=$serial_cursor --zone=$zone 2>&1 1>&3-) ;} 3>&1
  if [[ $cmd_err =~ $cursor_regex ]]
  then
    serial_cursor="${BASH_REMATCH[1]}"
  else
    echo "no cursor"
    break
  fi
done

echo -n "Waiting for VM shutdown"
until gcloud compute instances describe $instance --format='value("status")' --zone=$zone | grep -q "TERMINATED";
do 
  echo -n '.'
  sleep 3; 
done
echo
echo "deleting image maker VM"
gcloud compute instances delete $instance --zone=$zone --quiet & > /dev/null 2>&1
echo "Creating image"
gcloud compute images create ${image_name} --source-disk $instance --source-disk-zone=$zone
echo "Removing image maker disk"
gcloud compute disks delete $instance --zone=$zone --quiet & > /dev/null 2>&1
echo "Starting K3OS instance"
gcloud compute instances create k3os --zone=$zone --machine-type=n1-standard-1 --image=${image_name} --image-project=$project
