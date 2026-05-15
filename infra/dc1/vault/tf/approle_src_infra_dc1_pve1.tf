module "src_infra_dc1_pve1_approle" {
  source = "../../../../projects/tf_modules/vault_approle"
  name   = "src_infra_dc1_pve1"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets = vault_mount.secrets.path
  policies = [
    vault_policy.tf_token.name,
  ]
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}
