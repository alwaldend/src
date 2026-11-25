provider "vault" {
  skip_child_token = true
}

resource "random_password" "session_key" {
  length  = 64
  special = false
}

resource "vault_kv_secret_v2" "oidc" {
  mount = "secrets"
  name  = "alwaldend.com/vault1/approles/src_infra_mlflow/oidc"
  data_json = jsonencode({
    session_key = random_password.session_key.result
  })
}
