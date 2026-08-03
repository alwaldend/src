module "src_infra_forgejo_runner_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_infra_forgejo_runner"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets = vault_mount.secrets.path
  policies = [
    module.src_infra_forgejo_runner_ssh.policy,
  ]
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}

module "src_infra_forgejo_runner_ssh" {
  source          = "../../../projects/tf_modules/vault_ssh_server_role"
  backend         = vault_mount.ssh_servers.path
  name            = "src_infra_forgejo_runner_ssh"
  allowed_domains = "forgejo-runner.alwaldend.com"
}
