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
    vault_identity_entity.simeonwarrenbot.id,
  ]
  metadata = {
    comment = "Generic developer group"
  }
}
