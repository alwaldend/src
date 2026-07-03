locals {
  dns = jsondecode(file("${path.module}/../dnsconfig.json")).domains.default.records
}

module "vm_ha" {
  for_each = {
    host1 = { vmid = 600 },
  }
  source = ".././../../projects/tf_modules/pve_vm_qemu"
  name   = local.dns[each.key].A.name
  vmid   = each.value.vmid
  pool   = "src_infra_dc1_forgejo1"
  cores  = 2
  memory = 4096
  scsi0 = {
    size = "20G" # Boot
  }
  scsi1 = {
    size = "40G" # Forgejo
  }
  scsi2 = {
    size = "5G" # Traefik
  }
  ip   = "${local.dns[each.key].A.address}/24"
  tags = ["forgejo"]
}
