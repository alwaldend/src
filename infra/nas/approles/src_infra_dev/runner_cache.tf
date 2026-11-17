module "runner_cache" {
  source = "../../../../projects/tf_modules/truenas_dataset"
  parent = module.namespace.name
  config = {
    name   = "runner_cache"
    owner  = "src_infra_dev"
    access = { src_infra_xcp_ng = "MODIFY" }
    copies = 1
    quota  = 1024 * 1024 * 1024 * 1024
    smb = {
      name  = "src_infra_dev_runner_cache"
      hosts = ["192.168.1.213"]
    }
  }
}
