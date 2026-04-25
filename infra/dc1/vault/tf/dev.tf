resource "vault_identity_group" "dev" {
  name     = "dev"
  type     = "internal"
  policies = [vault_policy.dev.name]
  member_entity_ids = [data.vault_identity_entity.simeonwarren.id]

  metadata = {
    version = "2"
  }
}

resource "vault_policy" "dev" {
  name = "dev"

  policy = <<EOT
    path "secrets/data/entities/{{identity.entity.id}}/*" {
        capabilities = ["create", "update", "patch", "read", "delete"]
    }
    path "secrets/metadata/entities/{{identity.entity.id}}/*" {
      capabilities = ["list"]
    }
    path "auth/token/lookup-self" {
      capabilities = ["read"]
    }
    # Vault TF provider requires ability to create a child token
    path "auth/token/create" {
      capabilities = ["create", "update", "sudo"]
    }
EOT
}
