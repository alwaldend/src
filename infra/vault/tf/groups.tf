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
  metadata = {
    comment = "Generic developer group"
  }
}

resource "vault_identity_group" "sre" {
  name     = "sre"
  type     = "internal"
  policies = []
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  member_group_ids = [
    vault_identity_group.dev.id,
    vault_identity_group.ansible.id,
  ]
  metadata = {
    comment = "Generic admin group"
  }
}

resource "vault_identity_group" "ansible" {
  name = "ansible"
  type = "internal"
  policies = [
    vault_policy.ssh_clients_sign_admins.name,
  ]
  member_group_ids = [
    module.src_infra_dc1_forgejo1_approle.group_id,
    module.src_infra_harbor_approle.group_id,
    module.src_infra_flux_approle.group_id,
  ]
  metadata = {
    comment = "Group with access to ansible"
  }
}

resource "vault_identity_group" "users" {
  name = "users"
  type = "internal"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  metadata = {
    comment = "Users"
  }
}
