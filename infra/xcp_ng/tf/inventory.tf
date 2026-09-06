variable "resource_set_inventory" {
  description = "Pool-scoped names of the template, storage and network assigned to each AppRole. Unassigned sets remain empty."
  type = map(object({
    pool     = string
    template = string
    storage  = string
    network  = string
  }))
  default = {
    src_infra_dc1_forgejo1 = {
      pool     = "host1.xcp-ng.alwaldend.com"
      template = "fedora44-cloud.templates.xcp-ng.alwaldend.com"
      storage  = "Local storage"
      network  = "Pool-wide network 1"
    }
  }
}

data "xenorchestra_pool" "inventory" {
  for_each = toset(concat(
    [for inventory in var.resource_set_inventory : inventory.pool],
    [for vm in var.existing_vms : vm.pool],
  ))
  name_label = each.key
}

data "xenorchestra_template" "resource_sets" {
  for_each   = var.resource_set_inventory
  name_label = each.value.template
  pool_id    = data.xenorchestra_pool.inventory[each.value.pool].id
}

data "xenorchestra_sr" "resource_sets" {
  for_each   = var.resource_set_inventory
  name_label = each.value.storage
  pool_id    = data.xenorchestra_pool.inventory[each.value.pool].id
}

data "xenorchestra_network" "resource_sets" {
  for_each   = var.resource_set_inventory
  name_label = each.value.network
  pool_id    = data.xenorchestra_pool.inventory[each.value.pool].id
}
