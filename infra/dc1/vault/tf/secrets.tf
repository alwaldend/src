resource "vault_mount" "secrets" {
  path        = "secrets"
  type        = "kv"
  description = "Kv secrets"
  options = {
    version = "2"
  }
}
