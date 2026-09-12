module "src_infra_openhands" {
  source = "./approle_src_infra_openhands"
  name   = "src_infra_openhands"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets                  = vault_mount.secrets.path
  approle_backend          = vault_auth_backend.approle.path
  approle_backend_accessor = vault_auth_backend.approle.accessor
  ssh_backend              = vault_mount.ssh_servers.path
  ssh_allowed_domains      = "openhands.alwaldend.com"
  pki_backend              = module.pki_ica_servers.backend
  pki_allowed_domains = [
    "openhands.alwaldend.com",
    "canvas.openhands.alwaldend.com",
    "server.openhands.alwaldend.com",
    "automation.openhands.alwaldend.com",
  ]
}
