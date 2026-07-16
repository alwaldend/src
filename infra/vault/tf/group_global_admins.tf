resource "vault_identity_group" "global_admins" {
  name = "global_admins"
  type = "internal"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  member_group_ids = []
  metadata = {
    comment = "Global admins"
  }
}
