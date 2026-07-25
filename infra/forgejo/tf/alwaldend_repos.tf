resource "forgejo_repository" "alwaldend_src" {
  owner          = forgejo_organization.alwaldend.name
  name           = "src"
  description    = "Source code"
  website        = "https://alwaldend.com/"
  default_branch = "master"
  clone_addr     = "https://github.com/alwaldend/src.git"
}

resource "forgejo_collaborator" "alwaldend_src_flux" {
  repository_id = forgejo_repository.alwaldend_src.id
  user          = "src_infra_flux_git"
  permission    = "write"
}

resource "forgejo_branch_protection" "alwaldend_src_master" {
  branch_name           = "master"
  repository_id         = forgejo_repository.alwaldend_src.id
  enable_push           = true
  enable_push_whitelist = true
  push_whitelist_teams = [
    forgejo_team.alwaldend_admins.name,
  ]
  push_whitelist_usernames = [
    "src_infra_flux_git",
  ]
  enable_merge_whitelist = true
  merge_whitelist_teams = [
    forgejo_team.alwaldend_admins.name,
  ]
  enable_approvals_whitelist = true
  approvals_whitelist_teams = [
    forgejo_team.alwaldend_admins.name,
  ]
  require_signed_commits    = true
  block_on_rejected_reviews = true
  required_approvals        = 1
}

resource "forgejo_branch_protection" "alwaldend_src_releases" {
  branch_name           = "releases/*"
  repository_id         = forgejo_repository.alwaldend_src.id
  enable_push           = true
  enable_push_whitelist = true
  push_whitelist_teams = [
    forgejo_team.alwaldend_admins.name,
  ]
  enable_merge_whitelist = true
  merge_whitelist_teams = [
    forgejo_team.alwaldend_admins.name,
  ]
  enable_approvals_whitelist = true
  approvals_whitelist_teams = [
    forgejo_team.alwaldend_admins.name,
  ]
  require_signed_commits    = true
  block_on_rejected_reviews = true
  required_approvals        = 1
}
