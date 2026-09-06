data "vault_identity_group" "forgejo_access" {
  for_each   = toset(["admins", "package_writers", "src_writers"])
  group_name = "src_infra_dc1_forgejo1_${each.key}"

  lifecycle {
    postcondition {
      condition     = length(setsubtract(toset(self.member_entity_ids), local.vault_user_entity_ids)) == 0
      error_message = "Every Forgejo access-group member must belong to the Forgejo login group."
    }
    postcondition {
      condition     = each.key == "admins" || length(self.member_group_ids) == 0
      error_message = "Forgejo writer groups must contain direct entity memberships."
    }
  }
}

# The existing admin group includes its owning AppRole group. Resolve that
# group's direct entities, which include the service account and its operator.
data "vault_identity_group" "forgejo_admin_members" {
  for_each = toset(data.vault_identity_group.forgejo_access["admins"].member_group_ids)
  group_id = each.key

  lifecycle {
    postcondition {
      condition     = length(self.member_group_ids) == 0
      error_message = "Forgejo admin membership requires direct entities in the administrator group's children."
    }
    postcondition {
      condition     = length(setsubtract(toset(self.member_entity_ids), local.vault_user_entity_ids)) == 0
      error_message = "Every nested Forgejo administrator must belong to the Forgejo login group."
    }
  }
}

locals {
  alwaldend_src_writers = [
    for name, entity in local.vault_user_entities : forgejo_user.vault[name].login
    if contains(data.vault_identity_group.forgejo_access["src_writers"].member_entity_ids, entity.entity_id)
  ]
}
