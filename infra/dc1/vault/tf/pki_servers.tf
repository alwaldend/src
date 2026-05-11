module "pki_ica_servers" {
  depends_on   = [module.src_infra_dc1_vault_approle]
  source       = "../../../../projects/tf_modules/vault_pki_ica"
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

resource "vault_pki_secret_backend_role" "ica_servers_dc1_pve1" {
  backend            = module.pki_ica_servers.backend
  name               = "ica_servers_dc1_pve1"
  max_ttl            = local.month_in_seconds
  allow_ip_sans      = false
  key_type           = "rsa"
  key_bits           = 4096
  allow_any_name     = false
  allow_localhost    = false
  allowed_domains    = ["pve1.dc1.alwaldend.com"]
  allow_bare_domains = false
  allow_subdomains   = true
  server_flag        = true
  client_flag        = false
  no_store           = false
}
