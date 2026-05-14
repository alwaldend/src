resource "vault_mount" "ssh_clients" {
  path                      = "ssh/clients"
  type                      = "ssh"
  description               = "Ssh secrets: https://developer.hashicorp.com/vault/docs/secrets/ssh/signed-ssh-certificates"
  default_lease_ttl_seconds = local.year_in_seconds
  max_lease_ttl_seconds     = local.year_in_seconds
}

resource "vault_ssh_secret_backend_ca" "clients" {
  backend              = vault_mount.ssh_clients.path
  generate_signing_key = true
  key_type             = "ed25519"
}

resource "vault_ssh_secret_backend_role" "admins" {
  backend                 = vault_mount.ssh_clients.path
  name                    = "admins"
  key_type                = "ca"
  allow_user_certificates = true
  allowed_users_template  = true
  allowed_users           = "ansible,{{ identity.entity.name }}"
  default_user_template   = true
  default_user            = "{{ identity.entity.name }}"
  default_extensions = {
    permit-pty              = ""
    permit-agent-forwarding = ""
    permit-port-forwarding  = ""
  }
}

# Vault signs the key on each apply for some reason
# data "vault_ssh_secret_backend_sign" "simeonwarren" {
#   path       = vault_mount.ssh_clients.path
#   public_key = file("${path.module}/../../../../data/ssh/simeonwarren@alwaldend.com-2026-03-22-openpgp.pub")
#   name       = vault_ssh_secret_backend_role.admins.name
# }
#
# output "ssh_signed_key_simeonwarren" {
#   description = "Signed public key for simeonwarren"
#   value       = data.vault_ssh_secret_backend_sign.simeonwarren.signed_key
# }

resource "vault_ssh_secret_backend_role" "clients" {
  backend                 = vault_mount.ssh_clients.path
  name                    = "clients"
  key_type                = "ca"
  allow_user_certificates = true
  allowed_users_template  = true
  allowed_users           = "{{ identity.entity.name }}"
  default_user_template   = true
  default_user            = "{{ identity.entity.name }}"
  default_extensions = {
    permit-pty              = ""
    permit-agent-forwarding = ""
  }
}
