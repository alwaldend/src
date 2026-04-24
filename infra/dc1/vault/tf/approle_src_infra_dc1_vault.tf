resource "vault_identity_group" "src_infra_dc1_vault" {
  name     = "src_infra_dc1_vault"
  type     = "internal"
  policies = [vault_policy.src_infra_dc1_vault_approle.name]
  member_entity_ids = [vault_identity_entity.simeonwarren.id]

  metadata = {
    version = "2"
  }
}

resource "vault_approle_auth_backend_role" "src_infra_dc1_vault" {
  backend        = vault_auth_backend.approle.path
  role_name      = "src_infra_dc1_vault"
  role_id        = "src_infra_dc1_vault"
  secret_id_num_uses = 1
  secret_id_ttl = 3600
  token_policies = [vault_policy.src_infra_dc1_vault.name]
}

resource "vault_policy" "src_infra_dc1_vault_approle" {
  name = "src_infra_dc1_vault_approle"
  policy = <<EOT
    path "auth/${vault_auth_backend.approle.path}/role/${vault_approle_auth_backend_role.src_infra_dc1_vault.role_name}/secret-id" {
        capabilities = ["create"]
    }
EOT
}

resource "vault_policy" "src_infra_dc1_vault" {
  name = "src_infra_dc1_vault"
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
