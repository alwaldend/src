// Manage leases
// https://developer.hashicorp.com/vault/api-docs/system/leases
resource "vault_policy" "sys_leases_admin" {
  name   = "sys_leases_admin"
  policy = <<EOT
    path "sys/leases/*" {
      capabilities = ["create", "read", "update", "delete", "list", "sudo"]
    }
    path "sys/leases" {
      capabilities = [ "read" ]
    }
EOT
}
