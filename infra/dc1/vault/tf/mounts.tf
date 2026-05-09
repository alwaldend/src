resource "vault_mount" "secrets" {
  path        = "secrets"
  type        = "kv"
  description = "Kv secrets: https://developer.hashicorp.com/vault/docs/secrets/kv"
  options = {
    version = "2"
  }
}

resource "vault_mount" "ssh" {
  path        = "ssh/ssh1"
  type        = "ssh"
  description = "Ssh secrets: https://developer.hashicorp.com/vault/docs/secrets/ssh/signed-ssh-certificates"
}
