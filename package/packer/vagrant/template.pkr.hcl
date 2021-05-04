
variable "iso_checksum" {
  description = "SHA256 checksum of iso image"
  default     = "4DF22D4AEAFEEFD061BB0E33BF9E250E4991DD1853F3751B27B6792773EC1A2D"
  type        = string
}

variable "iso_url" {
  description = "URL to download k3OS iso image"
  default     = "https://github.com/rancher/k3os/releases/download/v0.20.4-k3s1r0/k3os-amd64.iso"
  type        = string
}

variable "box_version" {
  description = "k3OS box version"
  default     = "v0.20.4-k3s1r0"
  type        = string
}

variable "box_description" {
  description = "k3OS box description"
  default     = "k3OS is a Linux distribution designed to remove as much OS maintenance as possible in a Kubernetes cluster"
  type        = string
}

variable "password" {
  type    = string
  default = "rancher"
}

source "virtualbox-iso" "k3os" {
  boot_command = [
    "rancher", "<enter>",
    "sudo k3os install", "<enter>",
    "1", "<enter>",
    "y", "<enter>",
    "http://{{ .HTTPIP }}:{{ .HTTPPort }}/config.yml", "<enter>",
    "y", "<enter>"
  ]
  export_opts = [
    "--manifest",
    "--vsys", "0",
    "--description", "${var.box_description}",
    "--version", "${var.box_version}"
  ]
  boot_wait            = "40s"
  disk_size            = "8000"
  format               = "ova"
  guest_os_type        = "Linux_64"
  http_directory       = "."
  iso_checksum         = "sha256:${var.iso_checksum}"
  iso_url              = var.iso_url
  post_shutdown_delay  = "10s"
  shutdown_command     = "sudo poweroff"
  ssh_keypair_name     = ""
  ssh_private_key_file = "packer_rsa"
  ssh_timeout          = "1000s"
  ssh_username         = "rancher"
}

build {
  sources = ["source.virtualbox-iso.k3os"]

  post-processor "vagrant" {
    output = "k3os_{{.Provider}}.box"
  }
}
