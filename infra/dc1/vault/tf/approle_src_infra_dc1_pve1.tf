module "src_infra_dc1_pve1_approle" {
  source = "../../../../projects/tf_modules/vault_approle"
  name   = "src_infra_dc1_pve1"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets = vault_mount.secrets.path
  policies = [
    vault_policy.tf_token.name,
    module.src_infra_dc1_pve1_ssh.policy,
  ]
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}

module "src_infra_dc1_pve1_ssh" {
  source          = "../../../../projects/tf_modules/vault_ssh_server_role"
  backend         = vault_mount.ssh_servers.path
  name            = "src_infra_dc1_pve1_ssh"
  allowed_domains = "pve1.dc1.alwaldend.com"
}

module "src_infra_dc1_pve1_pki_server" {
  source          = "../../../../projects/tf_modules/vault_pki_server"
  backend         = module.pki_ica_servers.backend
  name            = "src_infra_dc1_pve1_pki_server"
  allowed_domains = ["pve1.dc1.alwaldend.com"]
}

resource "vault_identity_group" "src_infra_dc1_pve1_users" {
  name = "src_infra_dc1_pve1_users"
  type = "internal"
  member_group_ids = [
    vault_identity_group.dev.id,
    module.src_infra_dc1_pve1_approle.group_id,
  ]
}

resource "vault_identity_group" "src_infra_dc1_pve1_admins" {
  name = "src_infra_dc1_pve1_admins"
  type = "internal"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  member_group_ids = [
    module.src_infra_dc1_pve1_approle.group_id,
  ]
}
