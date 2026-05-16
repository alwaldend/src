resource "vault_policy" "auth_token_lookup_self" {
  name   = "auth_token_lookup_self"
  policy = <<EOT
    path "auth/token/lookup-self" {
      capabilities = ["read"]
    }
EOT
}

// Vault TF provider requires ability to create a child token
resource "vault_policy" "tf_token" {
  name   = "tf_token"
  policy = <<EOT
    path "auth/token/create" {
      capabilities = ["create", "update", "sudo"]
    }
EOT
}

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

// Allow OIDC auth: https://developer.hashicorp.com/vault/tutorials/auth-methods/oidc-identity-provider#configure-vault-authentication
resource "vault_policy" "identity_oidc_allow_auth" {
  name   = "identity_oidc_allow_auth"
  policy = <<EOT
    path "identity/oidc/provider/main/authorize" {
      capabilities = [ "read" ]
    }
EOT
}

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
