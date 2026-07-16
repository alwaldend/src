locals {
  dns = jsondecode(file("${path.module}/../dnsconfig.json")).domains.default.records
}

module "vms" {
  for_each = {
    host1 = { vmid = 1000 },
  }
  source = ".././../../projects/tf_modules/pve_vm_qemu"
  name   = "${local.dns[each.key].A.name}.alwaldend.com"
  vmid   = each.value.vmid
  pool   = "src_infra_threexui"
  cores  = 2
  memory = 4096
  scsi0 = {
    size = "20G" # Boot
  }
  scsi1 = {
    size = "20G" # Threexui
  }
  scsi2 = {
    size = "5G" # Traefik
  }
  ip   = "${local.dns[each.key].A.address}/24"
  tags = ["threexui"]
}
