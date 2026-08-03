locals {
  dns = jsondecode(file("${path.module}/../dnsconfig.json"))
  vms = {
    runner1 = {
      vmid = 1200
    }
  }
}

module "vm" {
  for_each = local.vms
  source   = "../../../projects/tf_modules/pve_vm_qemu"
  name     = "${local.dns.records[each.key].A.name}.alwaldend.com"
  vmid     = each.value.vmid
  pool     = "src_infra_forgejo_runner"
  cores    = 8
  memory   = 1024 * 16
  scsi0 = {
    size = "20G" # Boot
  }
  scsi1 = {
    size    = "300G" # Runner
    storage = "ceph-ec"
  }
  ip   = "${local.dns.records[each.key].A.address}/24"
  tags = ["forgejo-runner"]
}
