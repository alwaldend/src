resource "vault_identity_group" "truenas_user" {
  name = "truenas_user"
  type = "internal"
  member_entity_ids = setsubtract(
    setunion(
      vault_identity_group.approles.member_entity_ids,
      vault_identity_group.dns_approles.member_entity_ids,
    ),
    [data.vault_identity_entity.standalone_approles["src_infra_nas"].entity_id],
  )
}
