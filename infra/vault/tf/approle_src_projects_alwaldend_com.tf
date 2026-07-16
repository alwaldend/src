module "src_projects_alwaldend_com_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_projects_alwaldend_com"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets = vault_mount.secrets.path
  policies = [
    module.src_projects_alwaldend_com_ssh.policy,
  ]
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}

module "src_projects_alwaldend_com_ssh" {
  source          = "../../../projects/tf_modules/vault_ssh_server_role"
  backend         = vault_mount.ssh_servers.path
  name            = "src_projects_alwaldend_com_ssh"
  allowed_domains = "www.alwaldend.com"
}

module "src_projects_alwaldend_com_pki_server" {
  source                   = "../../../projects/tf_modules/vault_pki_server"
  backend                  = module.pki_ica_servers.backend
  name                     = "src_projects_alwaldend_com_pki_server"
  allowed_domains          = ["www.alwaldend.com"]
  eab_new_member_group_ids = [module.src_projects_alwaldend_com_approle.group_id]
  client_flag              = true
}

resource "vault_identity_group" "src_projects_alwaldend_com_users" {
  name              = "src_projects_alwaldend_com_users"
  type              = "internal"
  member_entity_ids = []
  member_group_ids = [
    vault_identity_group.src_projects_alwaldend_com_admins.id,
  ]
  metadata = {
    comment = "Users"
  }
}

resource "vault_identity_group" "src_projects_alwaldend_com_admins" {
  name = "src_projects_alwaldend_com_admins"
  type = "internal"
  member_group_ids = [
    module.src_projects_alwaldend_com_approle.group_id,
  ]
  metadata = {
    comment = "Admins"
  }
}
