resource "vault_mount" "ica1" {
  path                      = "pki/ica1"
  type                      = "pki"
  description               = "Intermediate Certificate Authority 1, signed externally: https://developer.hashicorp.com/vault/tutorials/pki/pki-engine-external-ca"
  default_lease_ttl_seconds = local.hour_in_seconds
  max_lease_ttl_seconds     = local.year_in_seconds * 3
}

resource "vault_pki_secret_backend_intermediate_cert_request" "ca1" {
  depends_on  = [module.src_infra_dc1_vault_approle, vault_mount.ica1]
  backend     = vault_mount.ica1.path
  type        = "internal"
  common_name = "Intermediate CA1"
  key_type    = "rsa"
  key_bits    = "4096"
}
resource "vault_pki_secret_backend_intermediate_set_signed" "ca1" {
  depends_on  = [vault_pki_secret_backend_intermediate_cert_request.ca1]
  backend     = vault_mount.ica1.path
  certificate = file("${path.module}/../../../../data/ssl/alwaldend.com/ica1.crt")
}

resource "vault_pki_secret_backend_config_issuers" "ca1_issuer" {
  depends_on                    = [vault_pki_secret_backend_intermediate_set_signed.ca1]
  backend                       = vault_mount.ica1.path
  default                       = vault_pki_secret_backend_intermediate_set_signed.ca1.imported_issuers[0]
  default_follows_latest_issuer = true
}
