module "src_infra_dc1_vault_approle" {
    source = "../../../../projects/tf_modules/vault_approle"
    name = "src_infra_dc1_vault"
    member_entity_ids = [data.vault_identity_entity.simeonwarren.id]
    secrets =  vault_mount.secrets.path
    policies = [vault_policy.tf_token.name]
    backend = vault_auth_backend.approle.path
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
        # Manage entities
        path "identity/*" {
           capabilities = ["create", "read", "update", "delete", "list"]
        }
        # Manage auth
        path "auth/*" {
           capabilities = ["create", "read", "update", "delete", "list"]
        }
EOT
}
