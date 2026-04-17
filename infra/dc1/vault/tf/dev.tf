resource "vault_mount" "dev" {
  path        = "dev"
  type        = "kv"
  description = "Secrets for developers"
  options = {
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
