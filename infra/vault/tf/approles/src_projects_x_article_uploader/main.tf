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
  name = "src_projects_x_article_uploader"
}

module "approle" {
  source                   = "../../../../../projects/tf_modules/vault_approle"
  name                     = local.name
  secrets                  = var.context.secrets
  backend                  = var.context.backend
  backend_accessor         = var.context.backend_accessor
  member_entity_ids        = var.context.member_entity_ids
  disable_yc_folder_policy = true
}

output "entity_id" {
  description = "AppRole entity id"
  value       = module.approle.entity_id
}
