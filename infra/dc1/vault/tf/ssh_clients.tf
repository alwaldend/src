# resource "vault_mount" "ssh_clients" {
#   path        = "ssh/admins"
#   type        = "ssh"
#   description = "Ssh secrets: https://developer.hashicorp.com/vault/docs/secrets/ssh/signed-ssh-certificates"
# }
#
# module "pki_ica_ssh_clients" {
#   depends_on   = [module.src_infra_dc1_vault_approle]
#   source       = "../../../../projects/tf_modules/vault_pki_ica"
#   name         = "ica_ssh_clients"
#   backend      = vault_mount.ica1.path
#   backend_cert = vault_pki_secret_backend_intermediate_set_signed.ca1.certificate
#   mount_path   = "pki/ica_ssh_clients"
# }
#
# resource "vault_ssh_secret_backend_ca" "clients" {
#     backend = vault_mount.ssh_clients.path
#     generate_signing_key = true
#     public_key = module.pki_ica_ssh_clients.certificate
#     private_key = module.pki_ica_ssh_clients.private_key
# }
#
# resource "vault_ssh_secret_backend_role" "admins" {
#     backend = vault_mount.ssh_clients.path
#     name                    = "admins"
#     key_type                = "ca"
#     allow_user_certificates = true
#     allowed_users_template = true
#     allowed_users = [
#         "ansible",
#         "{{ identity.entity.name }}",
#     ]
#     default_extensions = {
#         permit-pty = ""
#     }
# }
