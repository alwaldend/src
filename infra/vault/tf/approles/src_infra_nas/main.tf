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
  name = "src_infra_nas"
}

module "dns" {
  source  = "../../dns_access"
  name    = local.name
  secrets = var.context.secrets
  views   = ["dc1", "global"]
}

module "approle" {
  source                   = "../../../../../projects/tf_modules/vault_approle"
  name                     = local.name
  secrets                  = var.context.secrets
  backend                  = var.context.backend
  backend_accessor         = var.context.backend_accessor
  member_entity_ids        = var.context.member_entity_ids
  disable_yc_folder_policy = true
  policies                 = [module.dns.policy_name, vault_policy.nas_accounts.name]
}

resource "vault_policy" "nas_accounts" {
  name   = "src_infra_nas_accounts"
  policy = <<-EOT
    path "${var.context.secrets}/data/alwaldend.com/vault1/shared/read/approle/+/nas" {
      capabilities = ["create", "read", "update"]
    }
    path "${var.context.secrets}/metadata/alwaldend.com/vault1/shared/read/approle/+/nas" {
      capabilities = ["read"]
    }
  EOT
}

output "group_id" {
  description = "NAS provisioning AppRole group."
  value       = module.approle.group_id
}

output "entity_id" {
  description = "DNS owner AppRole entity id"
  value       = module.approle.entity_id
}
