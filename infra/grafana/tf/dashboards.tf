locals {
  dashboard_files = fileset("${path.module}/dashboards", "*/*/*.json")
  dashboard_configs = merge(
    { for file in local.dashboard_files : file => file("${path.module}/dashboards/${file}") },
    { "infra/victoria_metrics/monitoring-overview.json" = local.monitoring_overview },
  )
  project_paths = toset([for path in keys(local.dashboard_configs) : dirname(path)])
  namespaces    = toset([for path in local.project_paths : dirname(path)])
}

resource "grafana_folder" "namespace" {
  for_each                     = local.namespaces
  title                        = each.key
  uid                          = "repo-${substr(sha256(each.key), 0, 24)}"
  prevent_destroy_if_not_empty = true
}

resource "grafana_folder" "project" {
  for_each                     = local.project_paths
  title                        = basename(each.key)
  uid                          = "repo-${substr(sha256(each.key), 0, 24)}"
  parent_folder_uid            = grafana_folder.namespace[dirname(each.key)].uid
  prevent_destroy_if_not_empty = true
}

resource "grafana_dashboard" "repository" {
  for_each    = local.dashboard_configs
  folder      = grafana_folder.project[dirname(each.key)].uid
  config_json = each.value

  depends_on = [grafana_data_source.victoria_metrics]
}
