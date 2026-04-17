resource "vault_mount" "infra" {
  path        = "infra"
  type        = "kv"
  description = "Secrets for deployed infrastructure"
  options = {
    version = "2"
  }
}
