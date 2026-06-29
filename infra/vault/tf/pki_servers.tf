module "pki_ica_servers" {
  source       = "../../../projects/tf_modules/vault_pki_ica"
  name         = "ica_servers"
  backend      = vault_mount.ica1.path
  backend_cert = vault_pki_secret_backend_intermediate_set_signed.ca1.certificate
  mount_path   = "pki/ica_servers"
  // ACME: https://developer.hashicorp.com/vault/docs/secrets/pki/acme#enable-acme-support-on-a-vault-pki-mount
  // https://developer.hashicorp.com/vault/tutorials/pki/pki-acme-caddy#configure-acme
  allowed_response_headers = [
    "Last-Modified",
    "Location",
    "Replay-Nonce",
    "Link",
  ]
  passthrough_request_headers = ["If-Modified-Since"]
}

resource "vault_pki_secret_backend_config_cluster" "pki_ica_servers" {
  backend  = module.pki_ica_servers.backend
  path     = "${var.vault_url}/v1/${module.pki_ica_servers.backend}"
  aia_path = "${var.vault_url}/v1/${module.pki_ica_servers.backend}"
}

resource "vault_pki_secret_backend_config_acme" "pki_ica_servers" {
  backend                  = module.pki_ica_servers.backend
  enabled                  = true
  allowed_issuers          = ["*"]
  allowed_roles            = ["*"]
  allow_role_ext_key_usage = false
  default_directory_policy = "forbid"
  dns_resolver             = ""
  eab_policy               = "always-required"
}
