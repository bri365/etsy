data "vsphere_datacenter" "dc" {
  name = var.datacenter
}

data "vsphere_resource_pool" "pool" {
  name          = "${var.cluster}/Resources"
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_host" "hosts" {
  count         = length(var.hosts)
  name          = var.hosts[count.index]
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_datastore" "datastores" {
  count         = length(var.datastores)
  name          = var.datastores[count.index]
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_virtual_machine" "templates" {
  count         = length(var.templates)
  name          = var.templates[count.index]
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_network" "network" {
  name          = "VM Network"
  datacenter_id = data.vsphere_datacenter.dc.id
}
