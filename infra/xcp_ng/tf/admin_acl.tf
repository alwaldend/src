# Bootstrap terraform_data.oidc before the first full plan so these synchronized
# groups exist. Import pre-existing helper-managed ACLs by their XO acl.get IDs;
# adoption does not recreate or interrupt access.
locals {
  oidc_groups = jsondecode(data.external.oidc_identities.result.groups)
  pool_ids    = toset([for pool in data.xenorchestra_pools.all.pools : pool.id])
  existing_pool_admin_acl_ids = {
    for pool_id, acl_id in jsondecode(data.external.oidc_identities.result.pool_admin_acl_ids) : pool_id => acl_id
    if contains(local.pool_ids, pool_id)
  }
}

resource "xenorchestra_acl" "pool_administrator" {
  for_each = local.pool_ids
  subject  = lookup(local.oidc_groups, local.oidc_authorization.admin_group, "")
  object   = each.key
  action   = "admin"

  lifecycle {
    create_before_destroy = true

    precondition {
      condition     = contains(keys(local.oidc_groups), local.oidc_authorization.admin_group)
      error_message = "Bootstrap the XO OIDC plugin and synchronized administrator group before applying pool ACLs."
    }
  }
}

import {
  for_each = local.existing_pool_admin_acl_ids
  to       = xenorchestra_acl.pool_administrator[each.key]
  id       = each.value
}
