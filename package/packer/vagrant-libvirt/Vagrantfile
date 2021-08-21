# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.define "k3os"
  config.vm.box = "k3os"
  config.vm.guest = :linux
  config.vm.provider "libvirt" do |lv|
    lv.random_hostname = true
    lv.memory = 2048
    lv.cpus = 2
  end

  # Disable conflicting plugins
  config.vbguest.auto_update = false if Vagrant.has_plugin?("vagrant-vbguest") 

  # Disable default file syncing
  config.vm.synced_folder '.', '/vagrant', disabled: true
end
