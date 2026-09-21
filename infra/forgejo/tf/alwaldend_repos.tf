resource "forgejo_repository" "repositories" {
  for_each = module.repositories.forgejo_repositories

  owner           = forgejo_organization.alwaldend.name
  name            = each.value.name
  description     = each.value.description
  website         = each.value.homepage_url
  default_branch  = each.value.default_branch
  clone_addr      = each.value.config.clone_addr
  mirror          = try(each.value.config.mirror, false)
  mirror_interval = try(each.value.config.mirror_interval, null)
  has_actions     = try(each.value.config.has_actions, null)

  lifecycle {
    prevent_destroy = true
  }
}

moved {
  from = forgejo_repository.alwaldend_src
  to   = forgejo_repository.repositories["alwaldend/src"]
}

resource "forgejo_collaborator" "alwaldend_src_flux" {
  repository_id = forgejo_repository.repositories[local.src_repository_key].id
  user          = one(local.alwaldend_src_writers)
  permission    = "write"

  lifecycle {
    precondition {
      condition     = length(local.alwaldend_src_writers) == 1
      error_message = "The src repository automation writer group must contain exactly one discovered Forgejo user."
    }
  }
}

resource "forgejo_branch_protection" "alwaldend_src_master" {
  branch_name           = module.repositories.forgejo_repositories[local.src_repository_key].default_branch
  repository_id         = forgejo_repository.repositories[local.src_repository_key].id
  enable_push           = true
  enable_push_whitelist = true
  push_whitelist_teams = [
    forgejo_team.alwaldend_admins.name,
  ]
  push_whitelist_usernames = local.alwaldend_src_writers
  enable_merge_whitelist   = true
  merge_whitelist_teams = [
    forgejo_team.alwaldend_admins.name,
  ]
  enable_approvals_whitelist = true
  approvals_whitelist_teams = [
    forgejo_team.alwaldend_admins.name,
  ]
  require_signed_commits    = false
  block_on_rejected_reviews = true
  required_approvals        = 1
}

resource "forgejo_branch_protection" "alwaldend_src_releases" {
  branch_name           = "releases/*"
  repository_id         = forgejo_repository.repositories[local.src_repository_key].id
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
  require_signed_commits    = false
  block_on_rejected_reviews = true
  required_approvals        = 1
}
