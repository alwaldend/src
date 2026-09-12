module "dns_access" {
  source = "./dns_access"
  for_each = {
    src_infra_dc1_forgejo1     = ["dc1", "global"]
    src_infra_dc1_pve1         = ["dc1"]
    src_infra_dc1_vault        = ["dc1", "global"]
    src_infra_flux             = ["dc1"]
    src_infra_forgejo_runner   = ["dc1"]
    src_infra_harbor           = ["dc1"]
    src_infra_ingress          = ["dc1", "global"]
    src_infra_openhands        = ["dc1", "global"]
    src_infra_threexui         = ["dc1", "global"]
    src_infra_xcp_ng           = ["dc1"]
    src_projects_alwaldend_com = ["global"]
  }
  name    = each.key
  secrets = vault_mount.secrets.path
  views   = each.value
}
