data "vault_identity_entity" "standalone_approles" {
  for_each    = toset(["src_infra_nas", "src_infra_dev", "src_projects_al"])
  entity_name = each.key
}

data "vault_identity_group" "standalone_approles" {
  for_each   = toset(["src_infra_nas", "src_infra_dev"])
  group_name = each.key
}
