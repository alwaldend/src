module "vm_cloudinit_test" {
  source = ".././../../projects/tf_modules/pve_vm_qemu"
  name   = "${local.dns.domains.default.records.cloudinit_test.A.name}.alwaldend.com"
  vmid   = 100
  pool   = proxmox_pool.approles["src_infra_dc1_pve1"].poolid
  ip     = "${local.dns.domains.default.records.cloudinit_test.A.address}/24"
  scsi0 = {
    size = "5G"
  }
  tags = ["test"]
}
