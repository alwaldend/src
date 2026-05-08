locals {
  cacert_ca2 = format(
    "%s\n%s",
    vault_pki_secret_backend_root_sign_intermediate.ca2.certificate,
    local.cacert_ca1,
  )
}

resource "vault_pki_secret_backend_intermediate_cert_request" "ca2" {
  depends_on  = [module.src_infra_dc1_vault_approle, vault_mount.ica2]
  backend     = vault_mount.ica2.path
  type        = "internal"
  common_name = "Intermediate CA2"
  key_type    = "rsa"
  key_bits    = "4096"
}

resource "vault_pki_secret_backend_root_sign_intermediate" "ca2" {
  depends_on           = [vault_pki_secret_backend_intermediate_cert_request.ca2]
  backend              = vault_mount.ica1.path
  csr                  = vault_pki_secret_backend_intermediate_cert_request.ca2.csr
  common_name          = "Intermediate CA2"
  exclude_cn_from_sans = true
  max_path_length      = 1
  ttl                  = local.year_in_seconds
}

resource "vault_pki_secret_backend_intermediate_set_signed" "ca2" {
  depends_on  = [vault_pki_secret_backend_root_sign_intermediate.ca2]
  backend     = vault_mount.ica2.path
  certificate = local.cacert_ca2
}

resource "vault_pki_secret_backend_config_issuers" "ca2_issuer" {
  depends_on                    = [vault_pki_secret_backend_intermediate_set_signed.ca2]
  backend                       = vault_mount.ica2.path
  default                       = vault_pki_secret_backend_intermediate_set_signed.ca2.imported_issuers[0]
  default_follows_latest_issuer = true
}
