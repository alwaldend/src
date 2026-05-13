resource "vault_identity_oidc" "main" {
  issuer = var.vault_url
}

module "src_infra_dc1_pve1_provider" {
  source      = "../../../../projects/tf_modules/vault_oidc_provider"
  name        = "src_infra_dc1_pve1_provider"
  issuer_host = var.vault_host
  scopes_supported = [
    vault_identity_oidc_scope.user.name,
    vault_identity_oidc_scope.groups.name,
  ]
  group_ids = [
    vault_identity_group.src_infra_dc1_pve1_users.id,
  ]
  redirect_urls = [
    var.pve1_url,
  ]
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
