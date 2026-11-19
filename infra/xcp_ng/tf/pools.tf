data "vault_identity_group" "approles" {
  group_name = "approles"
}

data "vault_identity_entity" "approles_entities" {
  for_each  = toset(data.vault_identity_group.approles.member_entity_ids)
  entity_id = each.key
}

locals {
  entities_by_name = { for entity in data.vault_identity_entity.approles_entities : entity.entity_name => entity }
}

variable "resource_set_cpu_limit" {
  description = "CPU quota per AppRole resource set; XO requires at least one explicit quota."
  type        = number
  default     = 32

  validation {
    condition     = var.resource_set_cpu_limit > 0 && floor(var.resource_set_cpu_limit) == var.resource_set_cpu_limit
    error_message = "The CPU quota must be a positive integer."
  }
}

resource "xenorchestra_resource_set" "approles" {
  for_each = local.entities_by_name
  name     = each.key
  objects = contains(keys(var.resource_set_inventory), each.key) ? [
    data.xenorchestra_template.resource_sets[each.key].id,
    data.xenorchestra_sr.resource_sets[each.key].id,
    data.xenorchestra_network.resource_sets[each.key].id,
  ] : []
  subjects = contains(keys(local.oidc_users), each.value.entity_id) ? [local.oidc_users[each.value.entity_id]] : []

  limit {
    type     = "cpus"
    quantity = var.resource_set_cpu_limit
  }
}

output "approle_resource_sets" {
  description = "Xen Orchestra resource set IDs keyed by Vault AppRole entity name."
  value       = { for name, resource_set in xenorchestra_resource_set.approles : name => resource_set.id }
}
