resource "vault_identity_oidc" "main" {
  issuer = var.vault_url
}

resource "vault_identity_oidc_scope" "user" {
  name        = "user"
  template    = <<EOT
{
  "username": {{ identity.entity.metadata.username }},
  "nickname": {{ identity.entity.metadata.username }}
}
EOT
  description = "Username"
}

resource "vault_identity_oidc_scope" "email" {
  name        = "email"
  template    = <<EOT
{
  "email": {{ identity.entity.metadata.email }}
}
EOT
  description = "Email"
}

resource "vault_identity_oidc_scope" "ssh" {
  name        = "ssh"
  template    = <<EOT
{
  "sshpubkey": {{ identity.entity.metadata.sshpubkey }}
}
EOT
  description = "Public ssh key"
}

resource "vault_identity_oidc_scope" "groups" {
  name        = "groups"
  template    = <<EOT
{
  "groups": {{ identity.entity.groups.names }}
}
EOT
  description = "Groups"
}
