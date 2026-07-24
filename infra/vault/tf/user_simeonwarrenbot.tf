resource "vault_identity_entity" "simeonwarrenbot" {
  name = "simeonwarrenbot"
  metadata = {
    username  = "simeonwarrenbot"
    email     = "simeonwarrenbot@alwaldend.com"
    sshpubkey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIHTsu4hx9jG1wyClvbcp10hvsp6I7VISIxf3sbGbOtF3"
  }
}

resource "vault_identity_entity_alias" "userpass_simeonwarrenbot" {
  name           = vault_identity_entity.simeonwarrenbot.name
  mount_accessor = vault_auth_backend.userpass.accessor
  canonical_id   = vault_identity_entity.simeonwarrenbot.id
}
