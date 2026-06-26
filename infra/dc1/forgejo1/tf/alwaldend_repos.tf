resource "forgejo_repository" "alwaldend_src" {
  owner          = forgejo_organization.alwaldend.name
  name           = "src"
  description    = "Source code"
  website        = "https://alwaldend.com/"
  default_branch = "master"
  clone_addr     = "https://github.com/alwaldend/src.git"
}

locals {
  long_lived_branches = ["master", "releases/*"]
}

resource "forgejo_branch_protection" "alwaldend_src_long_lived" {
  for_each               = { for branch in local.long_lived_branches : branch => "" }
  branch_name            = each.key
  repository_id          = forgejo_repository.alwaldend_src.id
  enable_push            = false
  require_signed_commits = true
}
