# vSphere (vsphere_password and guest_password populated from ENV)
vsphere_server = "vcsa.rvidiot.io"
vsphere_user   = "administrator@vsphere.local"

datacenter     = "drv"
cluster        = "slide"
folder         = "etcd"

etcd_version = "3.4.13"
guest_user     = "brian"

# VM deployment targets
hosts      = ["esx1.rvidiot.io", "esx2.rvidiot.io", "esx3.rvidiot.io"]
datastores = ["esx1_nvme", "esx2_nvme", "esx3_nvme"]
templates  = ["ubu20-1", "ubu20-2", "ubu20-3"]
# templates  = ["ubu18", "ubu18", "ubu18"]

# Nodes to deploy (distributed across deployment targets)
etcds   = ["etcd1", "etcd2", "etcd3"]
# etcds   = ["etcd1",]

# Node configuration
vm_domain           = "rvidiot.io"
vm_base_address     = "192.168.50"
vm_starting_address = 181
vm_netmask          = 24
vm_gateway          = "192.168.50.1"
vm_dns_servers      = ["192.168.50.3",]
