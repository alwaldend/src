resource "vault_identity_group" "dev" {
  name = "dev"
  type = "internal"
  policies = [
    vault_policy.tf_token.name,
    vault_policy.auth_token_lookup_self.name,
    vault_policy.identity_oidc_allow_auth.name,
    vault_policy.ssh_clients_sign_clients.name,
  ]
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
}

resource "vault_identity_group" "sre" {
  name = "sre"
  type = "internal"
  policies = [
    vault_policy.tf_token.name,
    vault_policy.auth_token_lookup_self.name,
    vault_policy.identity_oidc_allow_auth.name,
    vault_policy.ssh_clients_sign_admins.name,
  ]
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
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
