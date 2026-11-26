locals {
  cluster_credential_ttl = 3 * 2629746 # Three months in seconds.
  oidc_scopes            = ["user", "groups", "email"]
}

data "vault_identity_group" "shared" {
  for_each   = toset(["dev", "approles"])
  group_name = each.key
}

module "src_infra_flux_pki_server" {
  source          = "../../../projects/tf_modules/vault_pki_server"
  backend         = "pki/ica_servers"
  name            = "src_infra_flux_pki_server"
  allowed_domains = ["flux.alwaldend.com"]
  eab_new_member_group_ids = [
    module.approle.group_id,
  ]
  signer_member_group_ids = [
    module.src_infra_flux_cluster_approle.group_id,
  ]
  client_flag = true
}

module "src_infra_flux_provider" {
  source           = "../../../projects/tf_modules/vault_oidc_provider"
  name             = "src_infra_flux_provider"
  issuer_host      = var.vault_host
  scopes_supported = local.oidc_scopes
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
    data.vault_identity_group.shared["dev"].group_id,
    data.vault_identity_group.shared["approles"].group_id,
    vault_identity_group.src_infra_flux_admins.id,
  ]
  metadata = {
    comment = "Users. Allowed to login with OIDC to flux"
  }
}

resource "vault_identity_group" "src_infra_flux_admins" {
  name              = "src_infra_flux_admins"
  type              = "internal"
  member_entity_ids = []
  member_group_ids = [
    module.approle.group_id,
  ]
  metadata = {
    comment = "flux admins"
  }
}

module "src_infra_flux_git_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_infra_flux_git"
  member_entity_ids = [
    module.approle.entity_id,
  ]
  secrets          = "secrets"
  sshpubkey        = "ecdsa-sha2-nistp384 AAAAE2VjZHNhLXNoYTItbmlzdHAzODQAAAAIbmlzdHAzODQAAABhBEpel7RvP2ELEdIJzcJlzK62FWCdHZHahL9G/+k8IdsD6OHcrw+m28xR4jnraqgMdC7ASWmWhFUX5lAjAyVNllfi4WUrtoJUUJDiZE3Ndclyyf8IwfgaeO75Xt0MW0fK7g=="
  backend          = data.vault_auth_backend.approle.path
  backend_accessor = data.vault_auth_backend.approle.accessor
}

module "src_infra_flux_cluster_approle" {
  source             = "../../../projects/tf_modules/vault_approle"
  name               = "src_infra_flux_cluster"
  secret_id_num_uses = 0
  secret_id_ttl      = local.cluster_credential_ttl
  token_ttl          = local.cluster_credential_ttl
  token_max_ttl      = local.cluster_credential_ttl
  member_entity_ids  = []
  member_group_ids = [
    module.approle.group_id,
  ]
  secrets          = "secrets"
  backend          = data.vault_auth_backend.approle.path
  backend_accessor = data.vault_auth_backend.approle.accessor
}

resource "vault_identity_group" "src_infra_flux_cluster_admins" {
  name              = "src_infra_flux_cluster_admins"
  type              = "internal"
  member_entity_ids = []
  member_group_ids = [
    module.src_infra_flux_cluster_approle.group_id,
  ]
  metadata = {
    comment = "Flux cluster admins"
  }
}

module "src_infra_flux_cluster_provider" {
  source           = "../../../projects/tf_modules/vault_oidc_provider"
  name             = "src_infra_flux_cluster_provider"
  issuer_host      = var.vault_host
  scopes_supported = local.oidc_scopes
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
  backend = "transit/default"
  name    = "src_infra_flux_cluster_sops"
  users_member_group_ids = [
    data.vault_identity_group.shared["dev"].group_id,
    data.vault_identity_group.shared["approles"].group_id,
  ]
  decryptors_member_group_ids = [
    module.src_infra_flux_cluster_approle.group_id,
  ]
}
