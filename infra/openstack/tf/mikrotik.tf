locals {
  inventory = yamldecode(
    file("${path.module}/../ansible/inventory.yaml")
  )

  master_hosts  = local.inventory.all.children.openstack_master.hosts
  compute_hosts = local.inventory.all.children.openstack_compute.hosts

  master_hostname  = keys(local.master_hosts)[0]
  compute_hostname = keys(local.compute_hosts)[0]

  master  = local.master_hosts[local.master_hostname]
  compute = local.compute_hosts[local.compute_hostname]

  leases = {
    master = {
      hostname    = local.master_hostname
      address     = local.master.ansible_host
      mac_address = local.master.openstack_mac_address
    }
    compute = {
      hostname    = local.compute_hostname
      address     = local.compute.ansible_host
      mac_address = local.compute.openstack_mac_address
    }
  }
}

resource "routeros_ip_dhcp_server_lease" "openstack" {
  for_each = local.leases

  address     = each.value.address
  mac_address = each.value.mac_address
  server      = "bridge1"
  lease_time  = "0s"
  comment     = "tf[infra/openstack/tf] ${each.value.hostname}"
}
