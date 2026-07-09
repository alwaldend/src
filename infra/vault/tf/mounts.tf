resource "vault_mount" "secrets" {
  path        = "secrets"
  type        = "kv"
  description = "Kv secrets: https://developer.hashicorp.com/vault/docs/secrets/kv"
  options = {
    version = "2"
  }
}

resource "vault_mount" "transit_default" {
  path                      = "transit/default"
  type                      = "transit"
  description               = "Transit: https://developer.hashicorp.com/vault/docs/secrets/transit"
  default_lease_ttl_seconds = 3600
  max_lease_ttl_seconds     = 86400
}
