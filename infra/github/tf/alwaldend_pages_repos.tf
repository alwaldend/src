locals {
  project_pages = jsondecode(file("${path.module}/project_pages.json"))
}

resource "github_repository" "project_landing" {
  for_each = local.project_pages

  name        = each.value
  description = "Landing page for ${each.key}"
  visibility  = "public"
  has_issues  = false
  has_wiki    = false
}

resource "github_repository_environment" "project_landing_pages" {
  for_each = local.project_pages

  repository  = github_repository.project_landing[each.key].name
  environment = "github-pages"
}
