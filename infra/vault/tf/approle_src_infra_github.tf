module "src_infra_github_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_infra_github"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets                  = vault_mount.secrets.path
  backend                  = vault_auth_backend.approle.path
  backend_accessor         = vault_auth_backend.approle.accessor
  disable_yc_folder_policy = true
}
