variable "context" {
  type = object({
    secrets           = string
    backend           = string
    backend_accessor  = string
    member_entity_ids = list(string)
  })
  description = "Vault mounts and operator identities supplied by the owning stage"
}

locals {
  name = "src_projects_rules_template"
}

module "dns" {
  source  = "../../dns_access"
  name    = local.name
  secrets = var.context.secrets
  views   = ["global"]
}

module "approle" {
  source                   = "../../../../../projects/tf_modules/vault_approle"
  name                     = local.name
  secrets                  = var.context.secrets
  backend                  = var.context.backend
  backend_accessor         = var.context.backend_accessor
  member_entity_ids        = var.context.member_entity_ids
  disable_yc_folder_policy = true
  policies                 = [module.dns.policy_name]
}

output "entity_id" {
  description = "DNS owner AppRole entity id"
  value       = module.approle.entity_id
}
