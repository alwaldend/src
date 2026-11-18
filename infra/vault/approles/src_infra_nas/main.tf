data "vault_auth_backend" "approle" {
  path = "approle"
}

data "vault_identity_entity" "operator" {
  entity_name = "simeonwarren"
}

locals {
  name = "src_infra_nas"
}

module "dns" {
  source  = "../../tf/dns_access"
  name    = local.name
  secrets = "secrets"
  views   = ["dc1", "global"]
}

module "approle" {
  source                   = "../../../../projects/tf_modules/vault_approle"
  name                     = local.name
  secrets                  = "secrets"
  backend                  = data.vault_auth_backend.approle.path
  backend_accessor         = data.vault_auth_backend.approle.accessor
  member_entity_ids        = [data.vault_identity_entity.operator.entity_id]
  disable_yc_folder_policy = true
  policies                 = [module.dns.policy_name, vault_policy.nas_accounts.name]
}

resource "vault_policy" "nas_accounts" {
  name   = "src_infra_nas_accounts"
  policy = <<-EOT
    path "secrets/data/alwaldend.com/vault1/shared/read/approle/+/nas" {
      capabilities = ["create", "read", "update"]
    }
    path "secrets/metadata/alwaldend.com/vault1/shared/read/approle/+/nas" {
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
