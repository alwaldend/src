resource "vault_mount" "pki" {
  path        = "pki"
  type        = "pki"
  description = "Pki mount"
  default_lease_ttl_seconds = 3600
  max_lease_ttl_seconds     = 86400
}

resource "vault_mount" "dev" {
    path        = "dev-secrets"
    type        = "kv"
    options     = {
        version = "2"
    }
}

resource "vault_policy" "dev" {
  name = "dev"

  policy = <<EOT
path "dev/+/creds" {
  capabilities = ["create", "update"]
}
path "dev/+/creds" {
  capabilities = ["read"]
}
## Vault TF provider requires ability to create a child token
path "auth/token/create" {
  capabilities = ["create", "update", "sudo"]
}
EOT
}
