module "src_infra_openstack_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_infra_openstack"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  policies = [
    module.src_infra_openstack_ssh.policy,
  ]
  secrets                  = vault_mount.secrets.path
  disable_yc_folder_policy = true
  backend                  = vault_auth_backend.approle.path
  backend_accessor         = vault_auth_backend.approle.accessor
}

module "src_infra_openstack_ssh" {
  source          = "../../../projects/tf_modules/vault_ssh_server_role"
  backend         = vault_mount.ssh_servers.path
  name            = "src_infra_openstack_ssh"
  allowed_domains = "openstack.alwaldend.com"
}
