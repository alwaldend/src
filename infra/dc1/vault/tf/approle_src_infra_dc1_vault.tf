module "src_infra_dc1_vault_approle" {
  source = "../../../../projects/tf_modules/vault_approle"
  name   = "src_infra_dc1_vault"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets = vault_mount.secrets.path
  policies = [
    vault_policy.tf_token.name,
    vault_policy.pki_admin.name,
    vault_policy.auth_admin.name,
    vault_policy.identity_admin.name,
    vault_policy.sys_leases_admin.name,
    vault_policy.sys_storage_admin.name,
    vault_policy.sys_mount_admin.name,
    vault_policy.sys_policies_acl_admin.name,
    vault_policy.ssh_admin.name,
  ]
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}
