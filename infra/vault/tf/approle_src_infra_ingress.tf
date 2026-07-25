module "src_infra_ingress_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_infra_ingress"
  member_group_ids = [
    vault_identity_group.global_admins.id,
  ]
  secrets          = vault_mount.secrets.path
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}

module "src_infra_ingress_ssh" {
  source          = "../../../projects/tf_modules/vault_ssh_server_role"
  backend         = vault_mount.ssh_servers.path
  name            = "src_infra_ingress_ssh"
  allowed_domains = "ingress.alwaldend.com"
}
