variable "vsphere_server" {
  default = "vcsa.example.com"
}

variable "vsphere_user" {
  default = "administrator@vsphere.local"
}

variable "vsphere_password" {
  description = "vCenter password from ENV"
}

variable "guest_user" {
  default = "ubuntu"
}

variable "guest_password" {
  description = "template VM password from ENV"
}

variable "datacenter" {
  default = "dc"
}

variable "cluster" {
  default = "cluster"
}

variable "folder" {
  default = "tf"
}

variable "hosts" {
  type = list
}

variable "datastores" {
  type = list
}

variable "masters" {
  type = list
}

variable "workers" {
  type = list
}

variable "templates" {
  type = list
}

variable "vm_domain" {
  default = "rvidiot.io"
}

variable "vm_base_address" {
  default = "192.168.50"
}

variable "vm_starting_address" {
  default = 181
}

variable "vm_netmask" {
  default = 24
}

variable "vm_gateway" {
  default = "192.168.50.1"
}

variable "vm_dns_servers" {
  type = list
}

variable "etcd_version" {
  default = "3.4.13"
}
