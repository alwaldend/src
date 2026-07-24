// Manage ACL
// https://developer.hashicorp.com/vault/api-docs/system/policies
resource "vault_policy" "sys_policies_acl_admin" {
  name   = "sys_policies_acl_admin"
  policy = <<EOT
    path "sys/policies/acl/*" {
      capabilities = ["create", "read", "update", "delete", "list", "sudo"]
    }
    path "sys/policies/acl" {
      capabilities = ["list"]
    }
EOT
}
