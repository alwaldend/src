terraform {
  required_providers {
    grafana = {
      source  = "grafana/grafana"
      version = "4.53.0"
    }
  }
}

locals {
  grafana_dns          = jsondecode(file("${path.module}/../dns/dnsconfig.json"))
  victoria_metrics_dns = jsondecode(file("${path.module}/../../victoria_metrics/dns/dnsconfig.json"))
}

provider "grafana" {
  url     = "https://${local.grafana_dns.records.local.A.name}"
  ca_cert = file("${path.module}/../../../data/ssl/alwaldend.com/root_ca.crt")
}

variable "victoria_metrics_username" {
  type        = string
  description = "HTTP Basic Auth username for the DC1 VictoriaMetrics endpoint."
  sensitive   = true
}

variable "victoria_metrics_password" {
  type        = string
  description = "HTTP Basic Auth password for the DC1 VictoriaMetrics endpoint."
  sensitive   = true
}

resource "grafana_data_source" "victoria_metrics" {
  name                = "VictoriaMetrics"
  uid                 = "victoria-metrics"
  type                = "prometheus"
  url                 = "https://${local.victoria_metrics_dns.records.local.A.name}.alwaldend.com"
  access_mode         = "proxy"
  basic_auth_enabled  = true
  basic_auth_username = var.victoria_metrics_username
  secure_json_data_encoded = jsonencode({
    basicAuthPassword = var.victoria_metrics_password
  })
}
