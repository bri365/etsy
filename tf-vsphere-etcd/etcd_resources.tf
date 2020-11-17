// Etcd nodes, distributed across hosts and datastores
resource "vsphere_virtual_machine" "etcd_nodes" {
  count            = length(var.etcds)

  name             = var.etcds[count.index]
  resource_pool_id = data.vsphere_resource_pool.pool.id
  host_system_id   = data.vsphere_host.hosts.*.id[count.index % length(var.hosts)]
  datastore_id     = data.vsphere_datastore.datastores.*.id[count.index % length(var.hosts)]
  folder           = var.folder

  num_cpus = 2
  memory   = 6000
  guest_id = "ubuntu64Guest"

  disk {
    label            = "disk0"
    size             = data.vsphere_virtual_machine.template_vm.disks.0.size
    eagerly_scrub    = data.vsphere_virtual_machine.template_vm.disks.0.eagerly_scrub
    thin_provisioned = data.vsphere_virtual_machine.template_vm.disks.0.thin_provisioned
  }

  network_interface {
    network_id   = data.vsphere_network.network.id
    # adapter_type = data.vsphere_virtual_machine.template_vm.network_interface_types[0]
  }

  clone {
    # template_uuid = data.vsphere_virtual_machine.template_20.id
    template_uuid = data.vsphere_virtual_machine.templates.*.id[count.index % length(var.hosts)]

    customize {
      linux_options {
        host_name = var.etcds[count.index]
        domain    = var.vm_domain
      }

      network_interface {
        ipv4_address = format("%s.%d", var.vm_base_address, var.vm_starting_address + count.index)
        ipv4_netmask = var.vm_netmask
      }

      ipv4_gateway    = var.vm_gateway
      dns_suffix_list = [var.vm_domain]
      dns_server_list = var.vm_dns_servers
    }
  }

  # upload etcd binaries
  provisioner "file" {
    source = "etcd-v${var.etcd_version}-linux-amd64.tar.gz"
    destination = "/home/${var.guest_user}/etcd-v${var.etcd_version}-linux-amd64.tar.gz"
  }

  # upload certificates
  provisioner "file" {
    source = "certs/etcd.tar.gz"
    destination = "/home/${var.guest_user}/etcd.tar.gz"
  }

  # upload systemd service file
  provisioner "file" {
    content = templatefile("${path.module}/etcd.service.tmpl",
      { name = var.etcds[count.index], ip = self.guest_ip_addresses.0 })
    destination = "/home/${var.guest_user}/etcd.service"
  }

  provisioner "remote-exec" {
    inline = [
        "sudo sed -i 's/127.0.1.1 ubu20/127.0.1.1 ${var.etcds[count.index]}/' /etc/hosts",
        "sudo groupadd --system etcd; sudo useradd -s /sbin/nologin --system -g etcd etcd",
        "sudo mkdir -p /var/lib/etcd; sudo mkdir -p /etc/etcd/certs",
        "sudo chown -R etcd:etcd /var/lib/etcd/",
        "sudo chmod -R 700 /var/lib/etcd/",
        "cd /home/${var.guest_user}; tar xf etcd-v${var.etcd_version}-linux-amd64.tar.gz",
        "cd /home/${var.guest_user}/etcd-v${var.etcd_version}-linux-amd64; sudo mv etcd etcdctl /usr/local/bin",
        "cd /home/${var.guest_user}; sudo tar xf etcd.tar.gz -C /etc/etcd/certs",
        "sudo chown -R etcd:etcd /etc/etcd/certs/",
        "sudo mv /home/${var.guest_user}/etcd.service /etc/systemd/system/etcd.service",
        "sudo systemctl daemon-reload; sudo systemctl enable etcd; sudo systemctl start etcd",
        "sleep 30; sudo systemctl status -l --no-pager etcd",
        "sudo etcdctl --cacert /etc/etcd/certs/ca.pem --cert /etc/etcd/certs/client.pem --key /etc/etcd/certs/client-key.pem member list",
      ]
  }

  connection {
    user = var.guest_user
    # password = var.guest_password
    host = self.guest_ip_addresses.0
    # private_key = file("perf_test_key.pem")
    agent = true
  }

}
