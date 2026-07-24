resource "vault_policy" "auth_token_lookup_self" {
  name   = "auth_token_lookup_self"
  policy = <<EOT
    path "auth/token/lookup-self" {
      capabilities = ["read"]
    }
EOT
}
