resource "vault_mount" "dev" {
  path        = "dev"
  type        = "kv"
  description = "Secrets for developers"
  options = {
    version = "2"
  }
}

resource "vault_identity_group" "dev" {
  name     = "dev"
  type     = "internal"
  policies = [vault_policy.dev.name]
  member_entity_ids = [vault_identity_entity.simeonwarren.id]

  metadata = {
    version = "2"
  }
}


resource "vault_policy" "dev" {
  name = "dev"

  policy = <<EOT
    path "dev/+/creds" {
      capabilities = ["create", "update"]
    }
    path "dev/+/creds" {
      capabilities = ["read"]
    }
    # Vault TF provider requires ability to create a child token
    path "auth/token/create" {
      capabilities = ["create", "update", "sudo"]
    }
EOT
}
