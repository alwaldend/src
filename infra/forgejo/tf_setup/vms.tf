module "vm_ha" {
  for_each = {
    host1 = { vmid = 600 },
  }
  source       = ".././../../projects/tf_modules/pve_vm_qemu"
  name         = local.dns.domains.default.records[each.key].A.name
  vmid         = each.value.vmid
  pool         = "src_infra_dc1_forgejo1"
  cores        = 4
  memory       = 8192
  storage_size = "60G"
  ipconfig0    = "ip=${local.dns.domains.default.records[each.key].A.address}/24,gw=192.168.10.2"
  tags         = ["forgejo"]
}
