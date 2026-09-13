# Adopt the protections created by the group defaults after verifying each
# repository import. Importing avoids unprotecting and recreating the rules.
resource "gitlab_branch_protection" "default_branches" {
  for_each = local.adopted_repositories

  project          = gitlab_project.repositories[each.key].id
  branch           = each.value.default_branch
  allow_force_push = false
  allowed_to_push = [{
    access_level = "maintainer"
  }]
  allowed_to_merge = [{
    access_level = "maintainer"
  }]

  lifecycle {
    prevent_destroy = true
  }
}

import {
  for_each = local.adopted_repositories

  to = gitlab_branch_protection.default_branches[each.key]
  id = "${each.value.config.id}:${each.value.default_branch}"
}
