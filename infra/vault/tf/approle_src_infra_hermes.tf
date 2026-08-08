module "src_infra_hermes_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_infra_hermes"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets = vault_mount.secrets.path
  policies = [
    module.src_infra_hermes_ssh.policy,
  ]
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}

module "src_infra_hermes_ssh" {
  source          = "../../../projects/tf_modules/vault_ssh_server_role"
  backend         = vault_mount.ssh_servers.path
  name            = "src_infra_hermes_ssh"
  allowed_domains = "hermes.alwaldend.com"
}

module "src_infra_hermes_pki_server" {
  source                   = "../../../projects/tf_modules/vault_pki_server"
  backend                  = module.pki_ica_servers.backend
  name                     = "src_infra_hermes_pki_server"
  allowed_domains          = ["hermes.alwaldend.com"]
  eab_new_member_group_ids = [module.src_infra_hermes_approle.group_id]
  client_flag              = true
}
