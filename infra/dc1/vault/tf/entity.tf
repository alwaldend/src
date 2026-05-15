resource "vault_identity_entity" "simeonwarren" {
  name = "simeonwarren"
  metadata = {
    username = "simeonwarren"
    email    = "simeonwarren@alwaldend.com"
  }
}

resource "vault_identity_entity_alias" "cert_simeonwarren" {
  name           = "simeonwarren.users.alwaldend.com"
  mount_accessor = vault_auth_backend.cert.accessor
  canonical_id   = vault_identity_entity.simeonwarren.id
}

resource "vault_cert_auth_backend_role" "cert_simeonwarren" {
  name                 = "simeonwarren"
  certificate          = file("${path.module}/../../../../data/ssl/alwaldend.com/simeonwarren.crt")
  backend              = vault_auth_backend.cert.path
  allowed_common_names = ["simeonwarren.users.alwaldend.com"]
  token_ttl            = local.day_in_seconds
  token_max_ttl        = local.day_in_seconds
}
