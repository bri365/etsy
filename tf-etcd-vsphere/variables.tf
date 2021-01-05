variable "vsphere_server" {
  description = "vCenter server"
}

variable "vsphere_user" {
  description = "vCenter user"
}

variable "vsphere_password" {
  description = "vCenter password"
}

variable "guest_user" {
  description = "template VM user"
}

variable "guest_password" {
  description = "template VM password"
}

variable "datacenter" {
  default = "dc"
}

variable "cluster" {
  default = "drv"
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

variable "etcds" {
  type = list
}

variable "templates" {
  type = list
}

variable "template_18" {
  default = "ubu18"
}

variable "template_20" {
  default = "ubu20"
}

variable "template_vm" {
  default = "ubu20"
}

variable "vm_domain" {
  default = "dc.tantrageek.com"
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
