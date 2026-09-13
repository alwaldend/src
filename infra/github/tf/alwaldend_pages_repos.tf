variable "bootstrap_projects" {
  description = "New projects whose pages branch must be published before enabling Pages. Leave empty after bootstrap."
  type        = set(string)
  default     = []
}

# Temporarily omit a live site's Pages configuration so GitHub reissues its
# custom-domain certificate, then restore it. This is a recovery action for an
# already-served site and must stay separate from first-creation bootstrap.
variable "certificate_reprovision_projects" {
  description = "Live projects whose Pages configuration is temporarily omitted to force certificate reissuance. Leave empty except during certificate recovery."
  type        = set(string)
  default     = []
}

locals {
  pages_disabled_projects = setunion(
    var.bootstrap_projects,
    var.certificate_reprovision_projects,
  )
}

resource "github_repository" "project_landing" {
  for_each = local.project_pages

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

  vulnerability_alerts = true

  dynamic "pages" {
    for_each = contains(local.pages_disabled_projects, each.key) ? [] : [each.value.config.pages]
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

    postcondition {
      condition     = self.repo_id == try(each.value.config.id, self.repo_id)
      error_message = "An adopted landing repository must retain its recorded numeric identity."
    }
  }
}

resource "github_repository_environment" "project_landing_pages" {
  for_each = local.project_pages

  repository  = github_repository.project_landing[each.key].name
  environment = "github-pages"

  deployment_branch_policy {
    protected_branches     = false
    custom_branch_policies = true
  }

  lifecycle {
    prevent_destroy = true
  }
}

# Keep Pages publication on its existing branch. The independent default
# branch starts with the same content, without renaming or deleting Pages.
resource "github_branch" "landing_default" {
  for_each = {
    for project, repository in local.project_pages : project => repository
    if !contains(var.bootstrap_projects, project)
  }

  repository    = github_repository.project_landing[each.key].name
  branch        = each.value.default_branch
  source_branch = each.value.config.pages.source.branch

  lifecycle {
    prevent_destroy = true

    precondition {
      condition     = each.value.default_branch != each.value.config.pages.source.branch
      error_message = "Landing repositories require separate default and Pages publication branches."
    }
  }
}
