resource "vault_identity_group" "dev" {
  name     = "dev"
  type     = "internal"
  policies = [vault_policy.dev.name, vault_policy.tf_token.name]
  member_entity_ids = [data.vault_identity_entity.simeonwarren.id]

  metadata = {
    version = "2"
  }
}

resource "vault_policy" "dev" {
  name = "dev"
  policy = <<EOT
    path "auth/token/lookup-self" {
      capabilities = ["read"]
    }
EOT
}
