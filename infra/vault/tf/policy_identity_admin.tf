// Manage identity
// https://developer.hashicorp.com/vault/api-docs/secret/identity
resource "vault_policy" "identity_admin" {
  name   = "identity_admin"
  policy = <<EOT
    path "identity/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    }
EOT
}
