data "vault_auth_backend" "approle" {
  path = "approle"
}

data "vault_identity_entity" "operator" {
  entity_name = "simeonwarren"
}

locals {
  name = "src_projects_al"
}

module "dns" {
  source  = "../../tf/dns_access"
  name    = local.name
  secrets = "secrets"
  views   = ["global"]
}

module "approle" {
  source                   = "../../../../projects/tf_modules/vault_approle"
  name                     = local.name
  secrets                  = "secrets"
  backend                  = data.vault_auth_backend.approle.path
  backend_accessor         = data.vault_auth_backend.approle.accessor
  member_entity_ids        = [data.vault_identity_entity.operator.entity_id]
  disable_yc_folder_policy = true
  policies                 = [module.dns.policy_name]
}

output "entity_id" {
  description = "DNS owner AppRole entity id"
  value       = module.approle.entity_id
}
