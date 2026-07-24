// Allow OIDC auth: https://developer.hashicorp.com/vault/tutorials/auth-methods/oidc-identity-provider#configure-vault-authentication
resource "vault_policy" "identity_oidc_allow_auth" {
  name   = "identity_oidc_allow_auth"
  policy = <<EOT
    path "identity/oidc/provider/main/authorize" {
      capabilities = [ "read" ]
    }
EOT
}
