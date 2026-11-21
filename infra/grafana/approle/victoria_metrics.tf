resource "vault_policy" "victoria_metrics_read" {
  name   = "src_infra_grafana_victoria_metrics_read"
  policy = <<-EOT
    path "secrets/data/alwaldend.com/vault1/approles/src_infra_victoria_metrics/web_auth" {
      capabilities = ["read"]
    }
  EOT
}
