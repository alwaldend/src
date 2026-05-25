module "vm_ha" {
  for_each = {
    host1 = { vmid = 500 },
    host2 = { vmid = 501 },
    host3 = { vmid = 502 }
  }
  source       = ".././../../../projects/tf_modules/pve_vm_qemu"
  name         = local.dns.domains.default.records[each.key].A.name
  vmid         = each.value.vmid
  pool         = "src_infra_dc1_consul1"
  cores        = 1
  memory       = 2560
  storage_size = "30G"
  ipconfig0    = "ip=${local.dns.domains.default.records[each.key].A.address}/24,gw=192.168.10.2"
  tags         = ["consul"]
}
