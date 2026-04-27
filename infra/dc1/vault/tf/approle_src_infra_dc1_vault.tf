module "src_infra_dc1_vault_approle" {
    source = "../../../../projects/tf_modules/vault_approle"
    name = "src_infra_dc1_vault"
    member_entity_ids = [data.vault_identity_entity.simeonwarren.id]
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
        # Vault TF provider requires ability to create a child token
        path "auth/token/create" {
          capabilities = ["create", "update", "sudo"]
        }
EOT
}

moved {
    from = vault_approle_auth_backend_role.src_infra_dc1_vault
    to = module.src_infra_dc1_vault_approle.vault_approle_auth_backend_role.role
}

moved {
    from = vault_policy.src_infra_dc1_vault
    to = module.src_infra_dc1_vault_approle.vault_policy.policy
}

moved {
    from = vault_policy.src_infra_dc1_vault_approle
    to = module.src_infra_dc1_vault_approle.vault_policy.policy_approle
}

moved {
    from = vault_identity_group.src_infra_dc1_vault
    to = module.src_infra_dc1_vault_approle.vault_identity_group.group
}
