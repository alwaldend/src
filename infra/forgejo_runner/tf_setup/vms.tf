locals {
  dns = jsondecode(file("${path.module}/../dnsconfig.json")).domains.default.records
}

module "vm" {
  source = "../../../projects/tf_modules/pve_vm_qemu"
  name   = "${local.dns.host1.A.name}.alwaldend.com"
  vmid   = 900
  pool   = "src_infra_forgejo_runner"
  cores  = 4
  memory = 8192
  scsi0 = {
    size = "40G"
  }
  ip   = "${local.dns.host1.A.address}/24"
  tags = ["forgejo-runner"]
}
