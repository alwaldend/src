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
  name = "src_infra_printers"
}

module "dns" {
  source  = "../../dns_access"
  name    = local.name
  secrets = var.context.secrets
  views   = ["dc1", "global"]
}

resource "vault_policy" "routeros_firewall" {
  name   = "${local.name}_routeros_firewall"
  policy = <<-EOT
    path "${var.context.secrets}/data/alwaldend.com/vault1/approles/src_infra_mikrotik/printers" {
      capabilities = ["read"]
    }
  EOT
}

module "approle" {
  source                   = "../../../../../projects/tf_modules/vault_approle"
  name                     = local.name
  secrets                  = var.context.secrets
  backend                  = var.context.backend
  backend_accessor         = var.context.backend_accessor
  member_entity_ids        = var.context.member_entity_ids
  disable_yc_folder_policy = true
  policies = [
    module.dns.policy_name,
    vault_policy.routeros_firewall.name,
  ]
}

output "entity_id" {
  description = "Printers owner AppRole entity id"
  value       = module.approle.entity_id
}
