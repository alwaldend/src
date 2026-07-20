locals {
  dns = jsondecode(file("${path.module}/../dnsconfig.json")).domains.default.records
}

module "vms" {
  for_each = {
    host1 = { vmid = 1100 },
  }
  source = ".././../../projects/tf_modules/pve_vm_qemu"
  name   = "${local.dns[each.key].A.name}.alwaldend.com"
  vmid   = each.value.vmid
  pool   = "user_simeonwarren"
  cores  = 4
  memory = 1024 * 8
  scsi0 = {
    size = "20G" # Boot
  }
  scsi1 = {
    size = "100G" # Opencode
  }
  scsi2 = {
    size = "5G" # Traefik
  }
  ip   = "${local.dns[each.key].A.address}/24"
  tags = ["opencode"]
}
