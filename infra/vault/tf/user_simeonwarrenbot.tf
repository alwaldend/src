resource "vault_identity_entity" "simeonwarrenbot" {
  name = "simeonwarrenbot"
  metadata = {
    username  = "simeonwarrenbot"
    email     = "simeonwarrenbot@alwaldend.com"
    sshpubkey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINPzqDtH8n4FV51GUzGV5TRQN3J+/GMUq63T4rJESTcV simeonwarrenbot@fedora"
  }
}

resource "vault_identity_entity_alias" "userpass_simeonwarrenbot" {
  name           = vault_identity_entity.simeonwarrenbot.name
  mount_accessor = vault_auth_backend.userpass.accessor
  canonical_id   = vault_identity_entity.simeonwarrenbot.id
}
