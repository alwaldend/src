// Manage mounts
// https://developer.hashicorp.com/vault/api-docs/system/mounts
// https://developer.hashicorp.com/vault/api-docs/system/remount
resource "vault_policy" "sys_mount_admin" {
  name   = "sys_mount_admin"
  policy = <<EOT
    path "sys/remount" {
      capabilities = ["create", "read", "update", "delete", "list", "sudo"]
    }
    path "sys/mounts/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    }
    path "sys/mounts" {
      capabilities = [ "read" ]
    }
EOT
}
