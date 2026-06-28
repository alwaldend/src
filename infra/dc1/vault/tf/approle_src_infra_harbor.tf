module "src_infra_harbor_approle" {
  source = "../../../../projects/tf_modules/vault_approle"
  name   = "src_infra_harbor"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets = vault_mount.secrets.path
  policies = [
    module.src_infra_harbor_ssh.policy,
  ]
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}

module "src_infra_harbor_ssh" {
  source          = "../../../../projects/tf_modules/vault_ssh_server_role"
  backend         = vault_mount.ssh_servers.path
  name            = "src_infra_harbor_ssh"
  allowed_domains = "harbor.alwaldend.com"
}

module "src_infra_harbor_pki_server" {
  source                   = "../../../../projects/tf_modules/vault_pki_server"
  backend                  = module.pki_ica_servers.backend
  name                     = "src_infra_harbor_pki_server"
  allowed_domains          = ["harbor.alwaldend.com"]
  eab_new_member_group_ids = [module.src_infra_harbor_approle.group_id]
  client_flag              = true
}

module "src_infra_harbor_provider" {
  source      = "../../../../projects/tf_modules/vault_oidc_provider"
  name        = "src_infra_harbor_provider"
  issuer_host = var.vault_host
  scopes_supported = [
    vault_identity_oidc_scope.user.name,
    vault_identity_oidc_scope.groups.name,
    vault_identity_oidc_scope.email.name,
  ]
  group_ids = [
    vault_identity_group.src_infra_harbor_users.id,
  ]
  allowed_read_clients_group_ids = [
    vault_identity_group.src_infra_harbor_admins.id,
  ]
  redirect_urls = [
    "${var.harbor_url}/c/oidc/callback",
  ]
}

resource "vault_identity_group" "src_infra_harbor_users" {
  name = "src_infra_harbor_users"
  type = "internal"
  member_group_ids = [
    vault_identity_group.dev.id,
    vault_identity_group.approles.id,
    vault_identity_group.src_infra_harbor_admins.id,
  ]
  metadata = {
    comment = "Users. Allowed to login with OIDC to Harbor"
  }
}

resource "vault_identity_group" "src_infra_harbor_admins" {
  name = "src_infra_harbor_admins"
  type = "internal"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  member_group_ids = [
    module.src_infra_harbor_approle.group_id,
  ]
  metadata = {
    comment = "Harbor admins"
  }
}
