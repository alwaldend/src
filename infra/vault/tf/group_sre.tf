resource "vault_identity_group" "sre" {
  name     = "sre"
  type     = "internal"
  policies = []
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  member_group_ids = [
  ]
  metadata = {
    comment = "Generic admin group"
  }
}
