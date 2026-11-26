locals {
  name = "src_infra_flux"
}

data "vault_auth_backend" "approle" {
  path = "approle"
}

data "vault_identity_entity" "operator" {
  entity_name = "simeonwarren"
}

module "approle" {
  source            = "../../../projects/tf_modules/vault_approle"
  name              = local.name
  member_entity_ids = [data.vault_identity_entity.operator.entity_id]
  secrets           = "secrets"
  policies = [
    module.dns.policy_name,
    module.ssh.policy,
  ]
  backend          = data.vault_auth_backend.approle.path
  backend_accessor = data.vault_auth_backend.approle.accessor
}

module "dns" {
  source  = "../../vault/tf/dns_access"
  name    = local.name
  secrets = "secrets"
  views   = ["dc1"]
}

module "ssh" {
  source          = "../../../projects/tf_modules/vault_ssh_server_role"
  backend         = "ssh/servers"
  name            = "${local.name}_ssh"
  allowed_domains = "flux.alwaldend.com"
}
