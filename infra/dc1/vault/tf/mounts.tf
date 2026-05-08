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

resource "vault_mount" "ica1" {
  path                      = "pki/ica1"
  type                      = "pki"
  description               = "Intermediate Certificate Authority 1, signed externally: https://developer.hashicorp.com/vault/tutorials/pki/pki-engine-external-ca"
  default_lease_ttl_seconds = local.hour_in_seconds
  max_lease_ttl_seconds     = local.year_in_seconds * 3
}

resource "vault_mount" "ica2" {
  path                      = "pki/ica2"
  type                      = "pki"
  description               = "Intermediate Certificate Authority 2, signed internally"
  default_lease_ttl_seconds = local.hour_in_seconds
  max_lease_ttl_seconds     = local.year_in_seconds
}
