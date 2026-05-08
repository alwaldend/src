locals {
  cacert_ca1 = file("${path.module}/cacerts/ca1.crt")
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
  certificate = local.cacert_ca1
}

resource "vault_pki_secret_backend_config_issuers" "ca1_issuer" {
  depends_on                    = [vault_pki_secret_backend_intermediate_set_signed.ca1]
  backend                       = vault_mount.ica1.path
  default                       = vault_pki_secret_backend_intermediate_set_signed.ca1.imported_issuers[0]
  default_follows_latest_issuer = true
}
