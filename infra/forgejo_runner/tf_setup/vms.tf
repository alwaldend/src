locals {
  runner_dns = jsondecode(file("${path.module}/../dnsconfig.json")).records

  runner_vms = {
    for name, vm in var.runners : name => {
      hostname = "${local.runner_dns[name].A.name}.alwaldend.com"
      address  = local.runner_dns[name].A.address
    }
  }
}

# Derive stable, locally administered MAC addresses from the owned hostname.
resource "xenorchestra_vm" "runner" {
  for_each = var.runners

  name_label       = local.runner_vms[each.key].hostname
  name_description = "Forgejo runner ${each.key} managed by infra/forgejo_runner"
  template         = data.xenorchestra_template.runner.id
  resource_set     = data.xenorchestra_resource_set.runner.id
  cpus             = each.value.cpus
  memory_max       = each.value.memory_gib * 1024 * 1024 * 1024
  memory_min       = each.value.memory_gib * 1024 * 1024 * 1024
  auto_poweron     = true
  power_state      = "Running"
  tags             = ["forgejo-runner", "src_infra_forgejo_runner", each.key]

  cloud_config = "#cloud-config\n${yamlencode(merge(
    jsondecode(file("${path.module}/../../cloud_init/xen_linux.json")),
    { hostname = split(".", local.runner_vms[each.key].hostname)[0], fqdn = local.runner_vms[each.key].hostname },
  ))}"

  cloud_network_config = yamlencode({
    version = 1
    config = [{
      type = "physical"
      name = var.xoa_guest_interface
      mac_address = join(":", concat(["02"], [
        for i in range(5) : substr(sha256(local.runner_vms[each.key].hostname), i * 2, 2)
      ]))
      subnets = [{
        type            = "static"
        address         = "${local.runner_vms[each.key].address}/24"
        gateway         = var.gateway
        dns_nameservers = var.dns_servers
      }]
    }]
  })

  network {
    network_id = data.xenorchestra_network.runner.id
    mac_address = join(":", concat(["02"], [
      for i in range(5) : substr(sha256(local.runner_vms[each.key].hostname), i * 2, 2)
    ]))
  }

  # Cloud-init grows the single root filesystem, including runner workspaces.
  disk {
    sr_id      = data.xenorchestra_sr.runner.id
    name_label = "${local.runner_vms[each.key].hostname} boot"
    size       = each.value.disk_gib * 1024 * 1024 * 1024
  }

  timeouts {
    create = "20m"
  }
}

output "vm_addresses" {
  description = "Static addresses of the created runner VMs keyed by runner name."
  value       = { for name, vm in local.runner_vms : name => vm.address }
}

output "vm_hostnames" {
  description = "Fully qualified hostnames of the created runner VMs keyed by runner name."
  value       = { for name, vm in local.runner_vms : name => vm.hostname }
}
