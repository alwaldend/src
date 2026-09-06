# XO creates provider-owned users during their first OIDC login. Bind exact
# identities: AppRole groups can contain other service entities.
# Run the OIDC configuration/bootstrap apply before this full configuration.
data "external" "oidc_identities" {
  program = [local.xo_config_binary, "oidc-identities"]
  query = {
    admin_group = local.oidc_authorization.admin_group
    url         = var.xoa_url
    insecure    = tostring(var.xoa_insecure)
    issuer      = trimsuffix(local.oidc_configuration.discoveryURL, "/.well-known/openid-configuration")
  }
}

locals {
  oidc_users = jsondecode(data.external.oidc_identities.result.users)
  existing_vm_grants = {
    for name, config in var.existing_vms : "${config.approle}/${one([
      for vm in data.xenorchestra_vms.existing[config.pool].vms : vm.id
      if vm.name_label == name && contains(vm.tags, config.approle)
      ])}" => {
      approle = config.approle
      vm_id = one([
        for vm in data.xenorchestra_vms.existing[config.pool].vms : vm.id
        if vm.name_label == name && contains(vm.tags, config.approle)
      ])
    }
  }
}

variable "existing_vms" {
  description = "Existing VM hostnames requiring direct owner ACLs; matches must also carry the owning AppRole tag."
  type = map(object({
    approle = string
    pool    = string
  }))
  default = {
    "host1.forgejo.alwaldend.com" = {
      approle = "src_infra_dc1_forgejo1"
      pool    = "host1.xcp-ng.alwaldend.com"
    }
  }
}

data "xenorchestra_vms" "existing" {
  for_each = toset([for vm in var.existing_vms : vm.pool])
  pool_id  = data.xenorchestra_pool.inventory[each.key].id

  lifecycle {
    postcondition {
      condition = alltrue([
        for name, config in var.existing_vms : length([
          for vm in self.vms : vm.id
          if vm.name_label == name && contains(vm.tags, config.approle)
        ]) == 1 if config.pool == each.key
      ])
      error_message = "Every existing VM must match exactly one hostname and owning AppRole tag in its named pool."
    }
  }
}

resource "xenorchestra_acl" "existing_vm_owner" {
  for_each = local.existing_vm_grants
  subject  = lookup(local.oidc_users, try(local.entities_by_name[each.value.approle].entity_id, ""), "")
  object   = each.value.vm_id
  action   = "admin"

  lifecycle {
    create_before_destroy = true

    precondition {
      condition     = contains(keys(local.entities_by_name), each.value.approle)
      error_message = "Every existing VM owner must be an entity in the Vault approles group."
    }
    precondition {
      condition     = contains(keys(local.oidc_users), try(local.entities_by_name[each.value.approle].entity_id, ""))
      error_message = "Log in to XO through Vault OIDC as the existing VM owner's AppRole before applying its ACL."
    }
  }
}

output "approles_pending_oidc_login" {
  description = "AppRoles with no synchronized XO OIDC user yet; their resource sets remain inaccessible until first login and the next apply."
  value       = sort([for name, entity in local.entities_by_name : name if !contains(keys(local.oidc_users), entity.entity_id)])
}
