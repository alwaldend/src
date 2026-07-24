// Ssh admin
// Docs: https://developer.hashicorp.com/vault/docs/secrets/ssh
// API: https://developer.hashicorp.com/vault/api-docs/secret/ssh
resource "vault_policy" "ssh_admin" {
  name   = "ssh_admin"
  policy = <<EOT
    path "ssh/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    }
    path "ssh/+/sign" {
      capabilities = ["read"]
    }
EOT
}
