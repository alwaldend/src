resource "vault_identity_entity" "simeonwarren" {
  name      = "simeonwarren"
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
