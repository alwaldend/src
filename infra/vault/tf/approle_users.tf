module "users_approle" {
  for_each          = local.users
  source            = "../../../projects/tf_modules/vault_approle"
  name              = "user_${each.value.name}"
  member_entity_ids = [each.value.id]
  secrets           = vault_mount.secrets.path
  backend           = vault_auth_backend.approle.path
  backend_accessor  = vault_auth_backend.approle.accessor
}
