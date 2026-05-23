module "vm_cloudinit_test" {
  source    = ".././../../../projects/tf_modules/pve_vm_qemu"
  name      = local.dns.domains.default.records.cloudinit_test.A.name
  vmid      = 100
  pool      = proxmox_pool.approles["src_infra_dc1_pve1"].poolid
  ipconfig0 = "ip=${local.dns.domains.default.records.cloudinit_test.A.address}/24,gw=192.168.10.2"
  tags      = ["test"]
}
