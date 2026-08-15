locals {
  inventory = yamldecode(
    file("${path.module}/../ansible/inventory.yaml")
  )
  master_vars = yamldecode(
    file("${path.module}/../ansible/group_vars/openstack_master.yaml")
  )
  compute_vars = yamldecode(
    file("${path.module}/../ansible/group_vars/openstack_compute.yaml")
  )

  master_hosts  = local.inventory.all.children.openstack_master.hosts
  compute_hosts = local.inventory.all.children.openstack_compute.hosts

  master_hostname  = keys(local.master_hosts)[0]
  compute_hostname = keys(local.compute_hosts)[0]

  leases = {
    master = {
      hostname    = local.master_hostname
      address     = local.master_vars.ansible_host
      mac_address = local.master_vars.openstack_mac_address
    }
    compute = {
      hostname    = local.compute_hostname
      address     = local.compute_vars.ansible_host
      mac_address = local.compute_vars.openstack_mac_address
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
