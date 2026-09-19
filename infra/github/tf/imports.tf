# Existing resource identities were observed before adoption. Entries with an
# observed GitHub ID must import successfully; they cannot be recreated.
import {
  to = github_organization_settings.alwaldend
  id = tostring(local.github_organization.github.id)
}

import {
  for_each = {
    for name, repository in local.repositories : name => repository
    if can(repository.config.id)
  }
  to = github_repository.other[each.key]
  id = each.value.name
}

import {
  for_each = {
    for name, repository in local.repositories : name => repository
    if can(repository.config.id)
  }
  to = github_branch_default.repository[each.key]
  id = each.value.name
}

import {
  for_each = local.organization_memberships
  to       = github_membership.organization[each.key]
  id       = "${local.github_organization_name}:${each.key}"
}

import {
  for_each = {
    for key, collaborator in local.developer_collaborators : key => collaborator
    if can(local.github_repositories[collaborator.repository].config.id)
  }
  to = github_repository_collaborator.developer[each.key]
  id = "${each.value.repository}:${each.value.username}"
}

import {
  to = github_repository_ruleset.src_main
  id = "src:6832533"
}

import {
  to = github_repository_ruleset.src_protected
  id = "src:19839736"
}
