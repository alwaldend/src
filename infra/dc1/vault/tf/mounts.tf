resource "vault_mount" "secrets" {
  path        = "secrets"
  type        = "kv"
  description = "Kv secrets: https://developer.hashicorp.com/vault/docs/secrets/kv"
  options = {
    version = "2"
  }
}

resource "vault_mount" "ssh" {
  path        = "ssh"
  type        = "ssh"
  description = "Ssh secrets: https://developer.hashicorp.com/vault/docs/secrets/ssh/signed-ssh-certificates"
}

resource "vault_mount" "pki" {
  path                      = "pki"
  type                      = "pki"
  description               = "Pki mount: https://developer.hashicorp.com/vault/docs/secrets/pki"
  default_lease_ttl_seconds = 3600
  max_lease_ttl_seconds     = 86400
}
