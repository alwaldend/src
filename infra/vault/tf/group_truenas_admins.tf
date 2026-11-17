resource "vault_identity_group" "truenas_admins" {
  name              = "truenas_admins"
  type              = "internal"
  member_entity_ids = [data.vault_identity_entity.standalone_approles["src_infra_nas"].entity_id]
}
