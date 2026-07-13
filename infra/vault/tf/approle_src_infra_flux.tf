module "src_infra_flux_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_infra_flux"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets = vault_mount.secrets.path
  policies = [
    module.src_infra_flux_ssh.policy,
  ]
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}

module "src_infra_flux_ssh" {
  source          = "../../../projects/tf_modules/vault_ssh_server_role"
  backend         = vault_mount.ssh_servers.path
  name            = "src_infra_flux_ssh"
  allowed_domains = "flux.alwaldend.com"
}

module "src_infra_flux_pki_server" {
  source          = "../../../projects/tf_modules/vault_pki_server"
  backend         = module.pki_ica_servers.backend
  name            = "src_infra_flux_pki_server"
  allowed_domains = ["flux.alwaldend.com"]
  eab_new_member_group_ids = [
    module.src_infra_flux_approle.group_id,
  ]
  signer_member_group_ids = [
    module.src_infra_flux_cluster_approle.group_id,
  ]
  client_flag = true
}

module "src_infra_flux_provider" {
  source      = "../../../projects/tf_modules/vault_oidc_provider"
  name        = "src_infra_flux_provider"
  issuer_host = var.vault_host
  scopes_supported = [
    vault_identity_oidc_scope.user.name,
    vault_identity_oidc_scope.groups.name,
    vault_identity_oidc_scope.email.name,
  ]
  group_ids = [
    vault_identity_group.src_infra_flux_users.id,
  ]
  allowed_read_clients_group_ids = [
    vault_identity_group.src_infra_flux_admins.id,
  ]
  redirect_urls = [
    "${var.flux_url}/oauth2/callback",
  ]
}

resource "vault_identity_group" "src_infra_flux_users" {
  name = "src_infra_flux_users"
  type = "internal"
  member_group_ids = [
    vault_identity_group.dev.id,
    vault_identity_group.approles.id,
    vault_identity_group.src_infra_flux_admins.id,
  ]
  metadata = {
    comment = "Users. Allowed to login with OIDC to flux"
  }
}

resource "vault_identity_group" "src_infra_flux_admins" {
  name = "src_infra_flux_admins"
  type = "internal"
  member_entity_ids = [
  ]
  member_group_ids = [
    module.src_infra_flux_approle.group_id,
  ]
  metadata = {
    comment = "flux admins"
  }
}

module "src_infra_flux_git_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_infra_flux_git"
  member_entity_ids = [
    module.src_infra_flux_approle.entity_id,
  ]
  secrets          = vault_mount.secrets.path
  sshpubkey        = "ecdsa-sha2-nistp384 AAAAE2VjZHNhLXNoYTItbmlzdHAzODQAAAAIbmlzdHAzODQAAABhBEpel7RvP2ELEdIJzcJlzK62FWCdHZHahL9G/+k8IdsD6OHcrw+m28xR4jnraqgMdC7ASWmWhFUX5lAjAyVNllfi4WUrtoJUUJDiZE3Ndclyyf8IwfgaeO75Xt0MW0fK7g=="
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}

module "src_infra_flux_cluster_approle" {
  source             = "../../../projects/tf_modules/vault_approle"
  name               = "src_infra_flux_cluster"
  secret_id_num_uses = 0
  secret_id_ttl      = local.month_in_seconds * 3
  token_ttl          = local.month_in_seconds * 3
  token_max_ttl      = local.month_in_seconds * 3
  member_entity_ids = [
  ]
  member_group_ids = [
    module.src_infra_flux_approle.group_id,
  ]
  secrets          = vault_mount.secrets.path
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}

resource "vault_identity_group" "src_infra_flux_cluster_admins" {
  name = "src_infra_flux_cluster_admins"
  type = "internal"
  member_entity_ids = [
  ]
  member_group_ids = [
    module.src_infra_flux_cluster_approle.group_id,
  ]
  metadata = {
    comment = "Flux cluster admins"
  }
}

module "src_infra_flux_cluster_provider" {
  source      = "../../../projects/tf_modules/vault_oidc_provider"
  name        = "src_infra_flux_cluster_provider"
  issuer_host = var.vault_host
  scopes_supported = [
    vault_identity_oidc_scope.user.name,
    vault_identity_oidc_scope.groups.name,
    vault_identity_oidc_scope.email.name,
  ]
  group_ids = [
    vault_identity_group.src_infra_flux_users.id,
  ]
  allowed_read_clients_group_ids = [
    vault_identity_group.src_infra_flux_admins.id,
  ]
  redirect_urls = [
    "https://unused",
  ]
}

module "src_infra_flux_cluster_sops" {
  source  = "../../../projects/tf_modules/vault_transit_key"
  backend = vault_mount.transit_default.path
  name    = "src_infra_flux_cluster_sops"
  users_member_group_ids = [
    resource.vault_identity_group.dev.id,
    resource.vault_identity_group.approles.id,
  ]
  decryptors_member_group_ids = [
    module.src_infra_flux_cluster_approle.group_id,
  ]
}
