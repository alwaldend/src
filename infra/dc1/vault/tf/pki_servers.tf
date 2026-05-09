module "pki_ica_servers" {
  depends_on   = [module.src_infra_dc1_vault_approle]
  source       = "../../../../projects/tf_modules/vault_pki_ica"
  name         = "ica_servers"
  backend      = vault_mount.ica1.path
  backend_cert = vault_pki_secret_backend_intermediate_set_signed.ca1.certificate
  mount_path   = "pki/ica_servers"
}
