resource "vault_identity_oidc" "main" {
  issuer = var.vault_url
}

resource "vault_identity_oidc_scope" "user" {
  name        = "user"
  template    = <<EOT
{
  "username": {{ identity.entity.name }}
}
EOT
  description = "The user scope provides claims using Vault identity entity metadata"
}

resource "vault_identity_oidc_scope" "groups" {
  name        = "groups"
  template    = <<EOT
{
  "groups": {{ identity.entity.groups.names }}
}
EOT
  description = "The groups scope provides the groups claim using Vault group membership"
}
