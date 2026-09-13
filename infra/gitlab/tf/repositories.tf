module "repositories" {
  source = "../../repos/tf"
}

locals {
  adopted_repositories = {
    for key, repository in module.repositories.gitlab_repositories :
    key => repository if try(repository.config.id, null) != null
  }

  organizations = {
    for name, organization in module.repositories.organizations :
    name => organization if try(organization.gitlab, null) != null
  }

  memberships = merge([
    for name, organization in local.organizations : {
      for username in setunion(toset(organization.admins), toset(organization.users)) :
      "${name}/${username}" => {
        organization = name
        username     = username
        access_level = contains(organization.admins, username) ? "owner" : "developer"
      }
    }
  ]...)

  existing_memberships = {
    for key, membership in local.memberships : key => membership
    if contains([
      for member in data.gitlab_group_membership.existing[membership.organization].members :
      member.username
    ], membership.username)
  }
}

data "gitlab_user" "members" {
  for_each = toset([for membership in local.memberships : membership.username])

  username = each.value
}

data "gitlab_group_membership" "existing" {
  for_each = local.organizations

  group_id  = each.value.gitlab.id
  inherited = false
}

resource "gitlab_group" "organizations" {
  for_each = local.organizations

  name                    = each.key
  path                    = each.key
  visibility_level        = each.value.gitlab.visibility
  project_creation_level  = "maintainer"
  subgroup_creation_level = "owner"

  default_branch_protection_defaults {
    allowed_to_push            = ["maintainer"]
    allowed_to_merge           = ["maintainer"]
    allow_force_push           = false
    developer_can_initial_push = false
  }

  lifecycle {
    prevent_destroy = true
  }
}

import {
  for_each = local.organizations

  to = gitlab_group.organizations[each.key]
  id = tostring(each.value.gitlab.id)
}

resource "gitlab_group_membership" "members" {
  for_each = local.memberships

  group_id     = gitlab_group.organizations[each.value.organization].id
  user_id      = data.gitlab_user.members[each.value.username].id
  access_level = each.value.access_level

  lifecycle {
    prevent_destroy = true
  }
}

import {
  for_each = local.existing_memberships

  to = gitlab_group_membership.members[each.key]
  id = "${local.organizations[each.value.organization].gitlab.id}:${data.gitlab_user.members[each.value.username].id}"
}

data "gitlab_project" "fork_sources" {
  for_each = toset([
    for repository in module.repositories.gitlab_repositories :
    repository.config.fork_from if repository.config.fork_from != null
  ])

  path_with_namespace = each.value
}

resource "gitlab_project" "repositories" {
  for_each = module.repositories.gitlab_repositories

  name                   = each.value.name
  path                   = each.value.name
  namespace_id           = gitlab_group.organizations[each.value.organization].id
  description            = each.value.description
  default_branch         = each.value.default_branch
  visibility_level       = each.value.visibility
  import_url             = each.value.config.import_url
  forked_from_project_id = each.value.config.fork_from == null ? null : tonumber(data.gitlab_project.fork_sources[each.value.config.fork_from].id)
  mr_default_target_self = each.value.config.fork_from == null ? null : false

  lifecycle {
    prevent_destroy = true

    postcondition {
      condition     = self.id == coalesce(try(tostring(each.value.config.id), null), self.id)
      error_message = "An adopted GitLab project must retain its recorded numeric identity."
    }
  }
}

import {
  for_each = local.adopted_repositories

  to = gitlab_project.repositories[each.key]
  id = tostring(each.value.config.id)
}
