// Manage transit
// https://developer.hashicorp.com/vault/docs/secrets/transit
resource "vault_policy" "transit_admin" {
  name   = "transit_admin"
  policy = <<EOT
    path "transit/*" {
      capabilities = ["create", "read", "update", "delete", "list", "sudo"]
    }
EOT
}
