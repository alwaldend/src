resource "vault_mount" "ssh_servers" {
  path                      = "ssh/servers"
  type                      = "ssh"
  description               = "Ssh secrets: https://developer.hashicorp.com/vault/docs/secrets/ssh/signed-ssh-certificates"
  default_lease_ttl_seconds = local.year_in_seconds
  max_lease_ttl_seconds     = local.year_in_seconds
}

resource "vault_ssh_secret_backend_ca" "servers" {
  backend              = vault_mount.ssh_servers.path
  generate_signing_key = true
  key_type             = "ed25519"
}

resource "vault_ssh_secret_backend_role" "servers" {
  backend                 = vault_mount.ssh_servers.path
  name                    = "servers"
  key_type                = "ca"
  allow_host_certificates = true
  allow_subdomains        = true
  allowed_domains         = ""
}
