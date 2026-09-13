resource "github_organization_settings" "alwaldend" {
  # Import is mandatory: the provider requires this argument even when billing
  # remains outside this module's ownership. Preserve its imported value.
  billing_email                 = ""
  default_repository_permission = local.github_organization.github.default_repository_permission

  lifecycle {
    prevent_destroy = true
    # Adopt organization identity and repository base access only. GitHub
    # remains authoritative for billing, profile, and unrelated defaults.
    ignore_changes = [
      billing_email,
      company,
      blog,
      email,
      twitter_username,
      location,
      name,
      description,
      has_organization_projects,
      has_repository_projects,
      members_can_create_repositories,
      members_can_create_public_repositories,
      members_can_create_private_repositories,
      members_can_create_internal_repositories,
      members_can_create_pages,
      members_can_create_public_pages,
      members_can_create_private_pages,
      members_can_fork_private_repositories,
      web_commit_signoff_required,
      advanced_security_enabled_for_new_repositories,
      dependabot_alerts_enabled_for_new_repositories,
      dependabot_security_updates_enabled_for_new_repositories,
      dependency_graph_enabled_for_new_repositories,
      secret_scanning_enabled_for_new_repositories,
      secret_scanning_push_protection_enabled_for_new_repositories,
    ]
  }
}

locals {
  organization_memberships = merge(
    { for username in local.github_organization.admins : username => "admin" },
    { for username in local.github_organization.users : username => "member" },
  )
  developer_collaborators = {
    for member in flatten([
      for name, repository in local.github_repositories : [
        for username in local.github_organization.users : {
          repository = name
          username   = username
        }
      ]
    ]) : "${member.repository}/${member.username}" => member
  }
}

resource "github_membership" "organization" {
  for_each = local.organization_memberships

  username = each.key
  role     = each.value

  lifecycle {
    prevent_destroy = true
  }
}

resource "github_repository_collaborator" "developer" {
  for_each = local.developer_collaborators

  repository = local.managed_repository_names[each.value.repository]
  username   = github_membership.organization[each.value.username].username
  permission = "push"

  lifecycle {
    prevent_destroy = true
  }
}
