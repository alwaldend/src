module "pki_ica_clients" {
  depends_on   = [module.src_infra_dc1_vault_approle]
  source       = "../../../projects/tf_modules/vault_pki_ica"
  name         = "ica_clients"
  backend      = vault_mount.ica1.path
  backend_cert = vault_pki_secret_backend_intermediate_set_signed.ca1.certificate
  mount_path   = "pki/ica_clients"
}

resource "vault_pki_secret_backend_role" "pki_ica_clients_users" {
  backend                  = module.pki_ica_clients.backend
  name                     = "pki_ica_clients_users"
  max_ttl                  = local.year_in_seconds
  allow_ip_sans            = true
  key_type                 = "rsa"
  key_bits                 = 4096
  key_usage                = ["DigitalSignature"]
  allow_any_name           = false
  allow_localhost          = false
  allowed_domains          = ["{{ identity.entity.metadata.username }}.users.alwaldend.com"]
  allowed_domains_template = true
  allow_bare_domains       = true
  allow_subdomains         = true
  server_flag              = false
  client_flag              = true
  no_store                 = false
}

resource "vault_policy" "pki_ica_clients_users_generator" {
  name   = "pki_ica_clients_users_generator"
  policy = <<EOT
    path "${vault_pki_secret_backend_role.pki_ica_clients_users.backend}/issue/${vault_pki_secret_backend_role.pki_ica_clients_users.name}" {
      capabilities = ["update"]
    }
EOT
}


resource "vault_identity_group" "pki_ica_clients_users_generator" {
  name = "pki_ica_clients_users_generate"
  type = "internal"
  policies = [
    vault_policy.pki_ica_clients_users_generator.name,
  ]
  member_group_ids = [
    vault_identity_group.users.id,
  ]
  metadata = {
    comment = "Group allowed to generate client certificates for themselves"
  }
}
