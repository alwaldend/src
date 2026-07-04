locals {
  dns = jsondecode(file("${path.module}/../dnsconfig.json"))
}

module "vms" {
  for_each = {
    host1 = { vmid = 800 },
  }
  source       = ".././../../projects/tf_modules/pve_vm_qemu"
  name         = local.dns.domains.default.records[each.key].A.name
  vmid         = each.value.vmid
  pool         = "src_infra_flux"
  cores        = 2
  memory       = 4096
  scsi0 = {
    size = "20G" # Boot
  }
  scsi1 = {
    size = "20G" # Local data
  }
  ip           = "${local.dns.domains.default.records[each.key].A.address}/24"
  tags         = ["flux"]
}
