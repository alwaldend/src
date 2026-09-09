locals {
  project_pages = jsondecode(file("${path.module}/project_pages.json"))
}

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

# The rules_skills project keeps the published rules-skill hostname and landing
# repository. Preserve the existing instances when the registry key changed
# from the legacy source name to the current one.
moved {
  from = github_repository.project_landing["rules_skill"]
  to   = github_repository.project_landing["rules_skills"]
}

moved {
  from = github_repository_environment.project_landing_pages["rules_skill"]
  to   = github_repository_environment.project_landing_pages["rules_skills"]
}

resource "github_repository" "project_landing" {
  for_each = local.project_pages

  name        = each.value.repository
  description = "Landing page for ${each.key}"
  visibility  = "public"
  has_issues  = false
  has_wiki    = false

  homepage_url = "https://${each.value.hostname}.alwaldend.com/"

  vulnerability_alerts = true

  dynamic "pages" {
    for_each = contains(local.pages_disabled_projects, each.key) ? [] : [true]
    content {
      build_type = "legacy"
      cname      = "${each.value.hostname}.alwaldend.com"
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
