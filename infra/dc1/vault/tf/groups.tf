resource "vault_identity_group" "dev" {
  name = "dev"
  type = "internal"
  policies = [
    vault_policy.tf_token.name,
    vault_policy.auth_token_lookup_self.name,
    vault_policy.identity_oidc_allow_auth.name,
  ]
  member_entity_ids = [
    data.vault_identity_entity.simeonwarren.id,
  ]
  metadata = {
    version = "2"
  }
}

resource "vault_identity_group" "src_infra_dc1_pve1_users" {
  name = "src_infra_dc1_pve1_users"
  type = "internal"
  member_group_ids = [
    vault_identity_group.dev.id,
  ]
  metadata = {
    version = "2"
  }
}

resource "vault_identity_group" "src_infra_dc1_pve1_admins" {
  name = "src_infra_dc1_pve1_admins"
  type = "internal"
  member_entity_ids = [
    data.vault_identity_entity.simeonwarren.id,
  ]
  metadata = {
    version = "2"
  }
}
