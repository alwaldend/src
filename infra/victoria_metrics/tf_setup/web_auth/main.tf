provider "vault" {
  skip_child_token = true
}

locals {
  username = "vm-${random_string.web_login.result}"
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

resource "vault_kv_secret_v2" "web_auth" {
  mount = "secrets"
  name  = "alwaldend.com/vault1/approles/src_infra_victoria_metrics/web_auth"
  data_json = jsonencode({
    username = local.username
    password = random_password.web_login.result
    htpasswd = "${local.username}:${random_password.web_login.bcrypt_hash}"
  })
}
