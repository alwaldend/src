resource "vault_policy" "web_login_read" {
  name   = "src_infra_victoria_metrics_web_login_read"
  policy = <<-EOT
    path "${module.approle.secret_mount}/data/${module.approle.secret_path}/web_auth" {
      capabilities = ["read"]
    }
  EOT
}

resource "vault_identity_group" "web_login_readers" {
  name             = "src_infra_victoria_metrics_web_login_readers"
  type             = "internal"
  member_group_ids = [data.vault_identity_group.operators.group_id]
  policies         = [vault_policy.web_login_read.name]
}
