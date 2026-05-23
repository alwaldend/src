locals {
  approles = {
    src_infra_dc1_consul1       = module.src_infra_dc1_consul1_approle
    src_infra_dc1_pve1          = module.src_infra_dc1_pve1_approle
    src_infra_dc1_vault         = module.src_infra_dc1_vault_approle
    src_infra_dc1_dns           = module.src_infra_dns_approle
    src_infra_yandex_cloud_org1 = module.src_infra_yandex_cloud_org1_approle
  }
}
