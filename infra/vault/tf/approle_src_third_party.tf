module "src_third_party_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_third_party"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets          = vault_mount.secrets.path
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}
