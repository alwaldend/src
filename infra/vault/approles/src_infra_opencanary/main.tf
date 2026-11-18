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
  policies                 = [module.ssh.policy]
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
