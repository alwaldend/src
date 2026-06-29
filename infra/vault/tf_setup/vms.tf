locals {
  dns = jsondecode(file("${path.module}/../dnsconfig.json"))
}

module "vm_ha" {
  for_each = {
    # host2 = { vmid = 400 },
    # host3 = { vmid = 401 }
  }
  source       = ".././../../projects/tf_modules/pve_vm_qemu"
  name         = local.dns.domains.default.records[each.key].A.name
  vmid         = each.value.vmid
  pool         = "src_infra_dc1_vault"
  cores        = 1
  memory       = 2560
  storage_size = "30G"
  ipconfig0    = "ip=${local.dns.domains.default.records[each.key].A.address}/24,gw=192.168.10.2"
  tags         = ["vault"]
}
