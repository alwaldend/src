module "src_infra_gitlab" {
  source = "./approles/src_infra_gitlab"
  context = {
    secrets           = vault_mount.secrets.path
    backend           = vault_auth_backend.approle.path
    backend_accessor  = vault_auth_backend.approle.accessor
    member_entity_ids = [vault_identity_entity.simeonwarren.id]
  }
}
