locals {
  project_pages = jsondecode(file("${path.module}/project_pages.json"))
}

variable "bootstrap_projects" {
  description = "New projects whose pages branch must be published before enabling Pages. Leave empty after bootstrap."
  type        = set(string)
  default     = []
}

resource "github_repository" "project_landing" {
  for_each = local.project_pages

  name        = each.value
  description = "Landing page for ${each.key}"
  visibility  = "public"
  has_issues  = false
  has_wiki    = false

  homepage_url = "https://${replace(each.key, "_", "-")}.alwaldend.com/"

  vulnerability_alerts = true

  dynamic "pages" {
    for_each = contains(var.bootstrap_projects, each.key) ? [] : [true]
    content {
      build_type = "legacy"
      cname      = "${replace(each.key, "_", "-")}.alwaldend.com"
      source {
        branch = "pages"
        path   = "/"
      }
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
}
