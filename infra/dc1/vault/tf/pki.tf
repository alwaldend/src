resource "vault_pki_secret_backend_intermediate_cert_request" "ca1" {
  backend     = vault_mount.pki.path
  type        = "internal"
  common_name = "Intermediate CA1"
  key_type    = "rsa"
  key_bits    = "4096"
}
