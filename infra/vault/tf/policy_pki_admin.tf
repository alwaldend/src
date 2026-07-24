// Manage pki
// https://developer.hashicorp.com/vault/api-docs/secret/pki
resource "vault_policy" "pki_admin" {
  name   = "pki_admin"
  policy = <<EOT
    path "pki/*" {
      capabilities = ["create", "read", "update", "delete", "list", "sudo"]
    }
EOT
}
