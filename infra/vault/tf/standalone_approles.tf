data "vault_identity_entity" "standalone_approles" {
  for_each    = toset(["src_infra_airflow", "src_infra_nas", "src_infra_dev", "src_infra_flux", "src_infra_flux_cluster", "src_infra_flux_git", "src_infra_grafana", "src_infra_mlflow", "src_infra_opencanary", "src_infra_victoria_metrics", "src_projects_al"])
  entity_name = each.key
}

data "vault_identity_group" "standalone_approles" {
  for_each   = toset(["src_infra_airflow", "src_infra_nas", "src_infra_dev", "src_infra_flux", "src_infra_flux_cluster", "src_infra_grafana", "src_infra_mlflow", "src_infra_opencanary", "src_infra_victoria_metrics"])
  group_name = each.key
}
