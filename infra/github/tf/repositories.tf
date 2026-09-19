module "repositories" {
  source = "../../repos/tf"
}

locals {
  github_organizations = {
    for name, organization in module.repositories.organizations : name => organization
    if can(organization.github)
  }
  # A provider instance manages one organization. Reject additional owners
  # until they have their own explicitly configured provider instance.
  github_organization_name = one(keys(local.github_organizations))
  github_organization      = local.github_organizations[local.github_organization_name]
  github_repositories = {
    for repository in module.repositories.github_repositories : repository.name => repository
  }
  # Landing repositories were retired with the project landing pages: the main
  # site now publishes every landing at /projects/<name>/, so the catalog no
  # longer carries a `landing_project` key and no dedicated Pages repository,
  # custom domain, or environment remains to manage.
  repositories = local.github_repositories
  managed_repository_names = {
    for name, repository in github_repository.other : name => repository.name
  }
}

resource "github_repository" "other" {
  for_each = local.repositories

  name                        = each.value.name
  description                 = each.value.description
  homepage_url                = each.value.homepage_url
  visibility                  = each.value.visibility
  has_issues                  = each.value.config.has_issues
  has_discussions             = try(each.value.config.has_discussions, false)
  has_projects                = each.value.config.has_projects
  has_wiki                    = each.value.config.has_wiki
  is_template                 = each.value.config.is_template
  archived                    = each.value.config.archived
  allow_merge_commit          = each.value.config.allow_merge_commit
  allow_squash_merge          = each.value.config.allow_squash_merge
  allow_rebase_merge          = each.value.config.allow_rebase_merge
  allow_auto_merge            = each.value.config.allow_auto_merge
  delete_branch_on_merge      = each.value.config.delete_branch_on_merge
  allow_update_branch         = each.value.config.allow_update_branch
  web_commit_signoff_required = each.value.config.web_commit_signoff_required
  squash_merge_commit_title   = each.value.config.squash_merge_commit_title
  squash_merge_commit_message = each.value.config.squash_merge_commit_message
  merge_commit_title          = each.value.config.merge_commit_title
  merge_commit_message        = each.value.config.merge_commit_message

  dynamic "pages" {
    for_each = try([each.value.config.pages], [])
    content {
      build_type = pages.value.build_type
      cname      = pages.value.cname
      dynamic "source" {
        for_each = pages.value.build_type == "legacy" ? [pages.value.source] : []
        content {
          branch = source.value.branch
          path   = source.value.path
        }
      }
    }
  }

  lifecycle {
    prevent_destroy = true
    # This adoption does not establish a new Dependabot alert policy.
    ignore_changes = [vulnerability_alerts]

    postcondition {
      condition     = self.repo_id == try(each.value.config.id, self.repo_id)
      error_message = "An adopted GitHub repository must retain its recorded numeric identity."
    }
  }
}

resource "github_branch_default" "repository" {
  for_each = local.repositories

  repository = local.managed_repository_names[each.key]
  branch     = each.value.default_branch

  lifecycle {
    prevent_destroy = true
  }
}
