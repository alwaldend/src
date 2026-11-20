data "vault_auth_backend" "approle" {
  path = "approle"
}

data "vault_identity_group" "operators" {
  group_name = "global_admins"
}

module "approle" {
  source                   = "../../../projects/tf_modules/vault_approle"
  name                     = "src_infra_victoria_metrics"
  member_group_ids         = [data.vault_identity_group.operators.group_id]
  policies                 = [module.dns.policy_name, module.ssh.policy]
  secrets                  = "secrets"
  disable_yc_folder_policy = true
  backend                  = data.vault_auth_backend.approle.path
  backend_accessor         = data.vault_auth_backend.approle.accessor
}

module "dns" {
  source  = "../../vault/tf/dns_access"
  name    = "src_infra_victoria_metrics"
  secrets = "secrets"
  views   = ["dc1", "global"]
}

module "pki_server" {
  source                   = "../../../projects/tf_modules/vault_pki_server"
  backend                  = "pki/ica_servers"
  name                     = "src_infra_victoria_metrics_pki_server"
  allowed_domains          = ["victoria-metrics.alwaldend.com"]
  allow_subdomains         = false
  eab_new_member_group_ids = [module.approle.group_id]
}

module "ssh" {
  source           = "../../../projects/tf_modules/vault_ssh_server_role"
  backend          = "ssh/servers"
  name             = "src_infra_victoria_metrics_ssh"
  allowed_domains  = "victoria-metrics,victoria-metrics.alwaldend.com"
  allow_subdomains = false
}
