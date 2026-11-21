locals {
  vmui_overview = jsondecode(file("${path.module}/../../victoria_metrics/ansible/files/monitoring-overview.json"))
  metric_units = {
    "req/s"     = "reqps"
    "samples/s" = "ops"
    "ms"        = "ms"
    "days"      = "d"
    "GiB"       = "gbytes"
  }
  overview_panels = flatten([
    for row_index, row in local.vmui_overview.rows : concat(
      [{
        title     = row.title
        type      = "row"
        collapsed = false
        panels    = []
        gridPos   = { x = 0, y = row_index * 9, w = 24, h = 1 }
      }],
      [for panel_index, panel in row.panels : {
        title       = panel.title
        description = try(panel.description, "")
        type        = "timeseries"
        datasource  = { type = "prometheus", uid = "victoria-metrics" }
        fieldConfig = {
          defaults  = { unit = lookup(local.metric_units, try(panel.unit, ""), "short") }
          overrides = []
        }
        gridPos = {
          x = 2 * sum(concat([0], [for previous in slice(row.panels, 0, panel_index) : previous.width]))
          y = row_index * 9 + 1
          w = 2 * panel.width
          h = 8
        }
        targets = [for query_index, expression in panel.expr : {
          refId        = "Q${query_index + 1}"
          expr         = expression
          legendFormat = try(panel.alias[query_index], "")
        }]
      }],
    )
  ])
  monitoring_overview = jsonencode({
    uid           = "infra-victoria-metrics-overview"
    title         = local.vmui_overview.title
    schemaVersion = 38
    tags          = ["infra/victoria_metrics"]
    timezone      = "browser"
    time          = { from = "now-1h", to = "now" }
    panels        = [for index, panel in local.overview_panels : merge(panel, { id = index + 1 })]
  })
}
