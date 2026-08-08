locals {
  dns = jsondecode(file("${path.module}/../dnsconfig.json")).domains.default.records
}

module "vm_ha" {
  for_each = {
    host1 = { vmid = 900 },
  }
  source = "../../../projects/tf_modules/pve_vm_qemu"
  name   = "${local.dns[each.key].A.name}.alwaldend.com"
  vmid   = each.value.vmid
  pool   = "src_infra_hermes"
  cores  = 4
  memory = 8192
  scsi0 = {
    size = "20G"
  }
  scsi1 = {
    size = "80G"
  }
  scsi2 = {
    size = "5G"
  }
  ip   = "${local.dns[each.key].A.address}/24"
  tags = ["hermes"]
}
