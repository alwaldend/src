module "src_projects_vastai_terraform_provider_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_projects_vastai_terraform_provider"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets                  = vault_mount.secrets.path
  backend                  = vault_auth_backend.approle.path
  backend_accessor         = vault_auth_backend.approle.accessor
  disable_yc_folder_policy = true
}
