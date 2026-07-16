module "src_infra_threexui_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_infra_threexui"
  member_group_ids = [
    vault_identity_group.global_admins.id,
  ]
  secrets = vault_mount.secrets.path
  policies = [
    module.src_infra_threexui_ssh.policy,
  ]
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}

module "src_infra_threexui_ssh" {
  source          = "../../../projects/tf_modules/vault_ssh_server_role"
  backend         = vault_mount.ssh_servers.path
  name            = "src_infra_threexui_ssh"
  allowed_domains = "threexui.alwaldend.com"
}

module "src_infra_threexui_pki_server" {
  source                   = "../../../projects/tf_modules/vault_pki_server"
  backend                  = module.pki_ica_servers.backend
  name                     = "src_infra_threexui_pki_server"
  allowed_domains          = ["threexui.alwaldend.com"]
  eab_new_member_group_ids = [module.src_infra_threexui_approle.group_id]
  client_flag              = true
}

resource "vault_identity_group" "src_infra_threexui_users" {
  name = "src_infra_threexui_users"
  type = "internal"
  member_group_ids = [
    vault_identity_group.src_infra_threexui_admins.id,
  ]
  metadata = {
    comment = "3x-ui Users"
  }
}

resource "vault_identity_group" "src_infra_threexui_admins" {
  name = "src_infra_threexui_admins"
  type = "internal"
  member_group_ids = [
    module.src_infra_threexui_approle.group_id,
  ]
  metadata = {
    comment = "3x-ui Admins"
  }
}
