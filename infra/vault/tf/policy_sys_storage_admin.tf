// Manage storage
// https://developer.hashicorp.com/vault/api-docs/system/storage
resource "vault_policy" "sys_storage_admin" {
  name   = "sys_storage_admin"
  policy = <<EOT
    path "sys/storage/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    }
    path "sys/storage" {
      capabilities = [ "read" ]
    }
EOT
}
