locals {
  dns    = jsondecode(file("${path.module}/../dnsconfig.json")).records
  server = var.vms["server"]

  vms = {
    for name, vm in var.vms : name => {
      hostname = "${local.dns[vm.dns_key].A.name}.alwaldend.com"
      address  = local.dns[vm.dns_key].A.address
    }
  }
}

# Every host signs its own host key with the shared OpenHands SSH role, so the
# MAC address is derived from the fully qualified hostname to stay stable and
# unique per component without a separate registry.
resource "xenorchestra_vm" "openhands" {
  for_each = var.vms

  name_label       = local.vms[each.key].hostname
  name_description = each.value.description
  template         = data.xenorchestra_template.openhands.id
  resource_set     = data.xenorchestra_resource_set.openhands.id
  cpus             = each.value.cpus
  memory_max       = each.value.memory_bytes
  memory_min       = each.value.memory_bytes
  auto_poweron     = true
  power_state      = "Running"
  tags             = each.value.tags

  cloud_config = "#cloud-config\n${yamlencode(merge(
    jsondecode(file("${path.module}/../../cloud_init/xen_linux.json")),
    { hostname = split(".", local.vms[each.key].hostname)[0], fqdn = local.vms[each.key].hostname },
  ))}"

  cloud_network_config = yamlencode({
    version = 1
    config = [{
      type = "physical"
      name = var.xoa_guest_interface
      mac_address = join(":", concat(["02"], [
        for i in range(5) : substr(sha256(local.vms[each.key].hostname), i * 2, 2)
      ]))
      subnets = [{
        type            = "static"
        address         = "${local.vms[each.key].address}/24"
        gateway         = var.gateway
        dns_nameservers = var.dns_servers
      }]
    }]
  })

  network {
    network_id = data.xenorchestra_network.openhands.id
    mac_address = join(":", concat(["02"], [
      for i in range(5) : substr(sha256(local.vms[each.key].hostname), i * 2, 2)
    ]))
  }

  # A single-disk cloud template maps its boot disk to xvda. The agent server
  # carries the conversations, workspaces and bash events for every backend
  # conversation, so it gets a larger root disk than the other components.
  disk {
    sr_id      = data.xenorchestra_sr.openhands.id
    name_label = "${local.vms[each.key].hostname} boot"
    size       = each.value.disk_bytes
  }

  timeouts {
    create = "20m"
  }
}

output "vm_addresses" {
  description = "Static addresses of the created OpenHands VMs keyed by component."
  value       = { for name, vm in local.vms : name => vm.address }
}

output "vm_hostnames" {
  description = "Fully qualified hostnames of the created OpenHands VMs keyed by component."
  value       = { for name, vm in local.vms : name => vm.hostname }
}
