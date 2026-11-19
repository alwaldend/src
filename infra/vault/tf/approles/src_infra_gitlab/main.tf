variable "context" {
  type = object({
    secrets           = string
    backend           = string
    backend_accessor  = string
    member_entity_ids = list(string)
  })
  description = "Vault mounts and operator identities supplied by the owning stage"
}

module "approle" {
  source                   = "../../../../../projects/tf_modules/vault_approle"
  name                     = "src_infra_gitlab"
  secrets                  = var.context.secrets
  backend                  = var.context.backend
  backend_accessor         = var.context.backend_accessor
  member_entity_ids        = var.context.member_entity_ids
  disable_yc_folder_policy = true
}
