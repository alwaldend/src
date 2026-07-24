// Manage auth
// https://developer.hashicorp.com/vault/api-docs/system/auth
// https://developer.hashicorp.com/vault/api-docs/auth
resource "vault_policy" "auth_admin" {
  name   = "auth_admin"
  policy = <<EOT
    path "auth/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    }
    path "sys/auth/*" {
      capabilities = ["create", "read", "update", "delete", "list", "sudo"]
    }
    path "sys/auth" {
      capabilities = [ "read" ]
    }
EOT
}
