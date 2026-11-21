data "vault_identity_entity" "standalone_approles" {
  for_each    = toset(["src_infra_nas", "src_infra_dev", "src_infra_grafana", "src_infra_opencanary", "src_infra_victoria_metrics", "src_projects_al"])
  entity_name = each.key
}

data "vault_identity_group" "standalone_approles" {
  for_each   = toset(["src_infra_nas", "src_infra_dev", "src_infra_grafana", "src_infra_opencanary", "src_infra_victoria_metrics"])
  group_name = each.key
}
