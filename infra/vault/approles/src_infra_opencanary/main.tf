data "vault_auth_backend" "approle" {
  path = "approle"
}

data "vault_identity_group" "operators" {
  group_name = "global_admins"
}

module "approle" {
  source                   = "../../../../projects/tf_modules/vault_approle"
  name                     = "src_infra_opencanary"
  member_group_ids         = [data.vault_identity_group.operators.group_id]
  policies                 = [module.dns.policy_name, module.ssh.policy]
  secrets                  = "secrets"
  disable_yc_folder_policy = true
  backend                  = data.vault_auth_backend.approle.path
  backend_accessor         = data.vault_auth_backend.approle.accessor
}

module "ssh" {
  source          = "../../../../projects/tf_modules/vault_ssh_server_role"
  backend         = "ssh/servers"
  name            = "src_infra_opencanary_ssh"
  allowed_domains = "opencanary.alwaldend.com"
}

module "pki_server" {
  source                  = "../../../../projects/tf_modules/vault_pki_server"
  backend                 = "pki/ica_servers"
  name                    = "src_infra_opencanary_pki_server"
  allowed_domains         = ["opencanary.alwaldend.com"]
  allow_subdomains        = false
  signer_member_group_ids = [module.approle.group_id]
}

module "dns" {
  source  = "../../tf/dns_access"
  name    = "src_infra_opencanary"
  secrets = "secrets"
  views   = ["dc1", "global"]
}
