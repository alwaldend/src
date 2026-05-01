# Vault TF provider requires ability to create a child token
resource "vault_policy" "tf_token" {
  name = "tf_token"
  policy = <<EOT
    path "auth/token/create" {
      capabilities = ["create", "update", "sudo"]
    }
EOT
}
