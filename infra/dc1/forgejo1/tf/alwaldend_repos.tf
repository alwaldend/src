resource "forgejo_repository" "alwaldend_src" {
  owner          = forgejo_organization.alwaldend.name
  name           = "src"
  description    = "Source code"
  website        = "https://alwaldend.com/"
  default_branch = "master"
  clone_addr     = "https://github.com/alwaldend/src.git"
}

resource "forgejo_branch_protection" "alwaldend_src_master" {
  branch_name            = "master"
  repository_id          = forgejo_repository.alwaldend_src.id
  enable_push            = true
  require_signed_commits = true
}

resource "forgejo_branch_protection" "alwaldend_src_releases" {
  branch_name            = "master"
  repository_id          = forgejo_repository.alwaldend_src.id
  enable_push            = false
  require_signed_commits = true
}
