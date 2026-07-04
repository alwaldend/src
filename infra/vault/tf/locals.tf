locals {
  year_in_seconds  = 31536000
  day_in_seconds   = 86400
  hour_in_seconds  = 3600
  month_in_seconds = 2629746
  users = {
    simeonwarren = vault_identity_entity.simeonwarren
  }
  approles = {
    src_infra_dc1_consul1       = module.src_infra_dc1_consul1_approle
    src_infra_dc1_forgejo1      = module.src_infra_dc1_forgejo1_approle
    src_infra_dc1_pve1          = module.src_infra_dc1_pve1_approle
    src_infra_dc1_router1       = module.src_infra_dc1_router1_approle
    src_infra_dc1_vault         = module.src_infra_dc1_vault_approle
    src_infra_dc1_dns           = module.src_infra_dns_approle
    src_infra_yandex_cloud_org1 = module.src_infra_yandex_cloud_org1_approle
    src_infra_harbor            = module.src_infra_harbor_approle
    src_infra_flux              = module.src_infra_flux_approle
    src_infra_flux_git          = module.src_infra_flux_git_approle
    users_simeonwarren          = module.users_approle["simeonwarren"]
  }
}
