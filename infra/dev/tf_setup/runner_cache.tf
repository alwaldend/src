locals {
  runner_cache_disk = one([
    for disk in xenorchestra_vm.secure1.disk : disk
    if disk.name_label == "${local.hostname} cache"
  ])
}

resource "vault_kv_secret_v2" "runner_cache" {
  mount = "secrets"
  name  = "alwaldend.com/vault1/approles/src_infra_dev/runner_cache"
  data_json = jsonencode({
    device = "/dev/xvd${substr("abcdefghijklmnopqrstuvwxyz", tonumber(local.runner_cache_disk.position), 1)}"
    size   = local.runner_cache_disk.size
  })

  lifecycle {
    precondition {
      condition = (
        local.runner_cache_disk.attached &&
        tonumber(local.runner_cache_disk.position) >= 1 &&
        tonumber(local.runner_cache_disk.position) <= 25
      )
      error_message = "The runner cache must be attached at a non-boot Xen disk position from 1 through 25."
    }
  }
}
