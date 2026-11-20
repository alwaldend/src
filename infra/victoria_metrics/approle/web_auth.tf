locals {
  web_login_path    = "alwaldend.com/vault1/victoria_metrics/web_login"
  web_verifier_path = "alwaldend.com/vault1/approles/src_infra_victoria_metrics/web_auth"
  web_login_name    = "vm-${random_string.web_login.result}"
}

resource "random_string" "web_login" {
  length  = 12
  upper   = false
  special = false
}

resource "random_password" "web_login" {
  # The bcrypt verifier supports up to 72 bytes.
  length  = 72
  special = false
}

resource "vault_kv_secret_v2" "web_login" {
  mount = "secrets"
  name  = local.web_login_path
  data_json = jsonencode({
    username = local.web_login_name
    password = random_password.web_login.result
  })
}

resource "vault_kv_secret_v2" "web_verifier" {
  mount = "secrets"
  name  = local.web_verifier_path
  data_json = jsonencode({
    htpasswd = "${local.web_login_name}:${random_password.web_login.bcrypt_hash}"
  })
}

resource "vault_policy" "web_login_read" {
  name   = "src_infra_victoria_metrics_web_login_read"
  policy = <<-EOT
    path "secrets/data/${local.web_login_path}" {
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
