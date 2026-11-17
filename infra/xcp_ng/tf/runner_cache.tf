# Creation-only provider gap: no SR resource exists in xenorchestra 0.41.0.
# The helper preserves an owned matching SR and never deletes it on rollback.
resource "terraform_data" "runner_cache" {
  provisioner "local-exec" {
    interpreter = [local.xo_config_binary]
    command     = "ensure-cache-storage"
    environment = {
      XOA_URL            = var.xoa_url
      XOA_INSECURE       = tostring(var.xoa_insecure)
      XO_CACHE_HOST      = "host1.xcp-ng.alwaldend.com"
      XO_CACHE_SR        = "Runner cache"
      XO_CACHE_SMB_SHARE = "//nas.alwaldend.com/src_infra_dev_runner_cache"
    }
    quiet = true
  }

  lifecycle {
    prevent_destroy = true
  }
}

data "xenorchestra_sr" "runner_cache" {
  name_label = "Runner cache"
  pool_id    = data.xenorchestra_pool.inventory["host1.xcp-ng.alwaldend.com"].id
  depends_on = [terraform_data.runner_cache]
}
