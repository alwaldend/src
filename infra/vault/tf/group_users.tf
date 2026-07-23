resource "vault_identity_group" "users" {
  name = "users"
  type = "internal"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
    vault_identity_entity.simeonwarrenbot.id,
  ]
  metadata = {
    comment = "Users"
  }
}
