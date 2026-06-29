resource "vault_mount" "secrets" {
  path        = "secrets"
  type        = "kv"
  description = "Kv secrets: https://developer.hashicorp.com/vault/docs/secrets/kv"
  options = {
    version = "2"
  }
}
