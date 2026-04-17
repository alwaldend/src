resource "vault_mount" "sre" {
  path        = "sre"
  type        = "kv"
  description = "Secrets for SRE"
  options = {
    version = "2"
  }
}

resource "vault_policy" "sre_infra_dc1_vault" {
  name = "sre_infra_dc1_vault"

  policy = <<EOT
    # Manage namespaces at root namespace level
    path "sys/namespaces/*" {
       capabilities = ["create", "read", "update", "delete", "list", "sudo"]
    }
    # List namespaces at root namespace level
    path "sys/namespaces" {
       capabilities = ["list"]
    }

    # Manage policies at root namespace level
    path "sys/policies/acl/*" {
       capabilities = ["create", "read", "update", "delete", "list", "sudo"]
    }

    # List policies at root namespace level
    path "sys/policies/acl" {
      capabilities = ["list"]
    }

    # Enable and manage secrets engines at root namespace level
    path "sys/mounts/*" {
       capabilities = ["create", "read", "update", "delete", "list"]
    }

    # List available secrets engines at root namespace level
    path "sys/mounts" {
      capabilities = [ "read" ]
    }

    # Enable and manage auth methods at root namespace level
    path "sys/auth/*" {
       capabilities = ["create", "read", "update", "delete", "list"]
    }

    # List available auth methods at root namespace level
    path "sys/auth" {
      capabilities = [ "read" ]
    }
EOT
}
