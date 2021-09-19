
variable "box_description" {
  type    = string
  default = "k3OS is a Linux distribution designed to remove as much OS maintenance as possible in a Kubernetes cluster"
}

variable "box_version" {
  type    = string
  default = "v0.20.7-k3s1r0"
}

variable "iso_checksum" {
  type    = string
  default = "85a560585bc5520a793365d70e6ce984f3fb2ce5a43b31f0f7833dc347487e69"
}

variable "iso_url" {
  type    = string
  default = "https://github.com/rancher/k3os/releases/download/v0.20.7-k3s1r0/k3os-amd64.iso"
}

variable "password" {
  type    = string
  default = "rancher"
}

source "qemu" "k3os" {
  boot_command = [
    "rancher", "<enter>",
    "sudo k3os install", "<enter>",
    "1", "<enter>",
    "2", "<enter>", // avoid fd0
    "y", "<enter>",
    "http://{{ .HTTPIP }}:{{ .HTTPPort }}/config.yml", "<enter>",
    "y", "<enter>",
  ]
  boot_wait            = "60s"
  disk_size            = "8000"
  disk_interface       = "virtio"
  format               = "qcow2"
  http_directory       = "."
  iso_checksum         = "sha256:${var.iso_checksum}"
  iso_url              = "${var.iso_url}"
  net_device           = "virtio-net"
  shutdown_command     = "sudo poweroff"
  ssh_keypair_name     = ""
  ssh_private_key_file = "packer_rsa"
  ssh_timeout          = "1000s"
  ssh_username         = "rancher"
}

build {
  sources = ["source.qemu.k3os"]

  post-processor "vagrant" {
    output            = "k3os_{{.Provider}}.box"
    provider_override = "libvirt"
  }
}
