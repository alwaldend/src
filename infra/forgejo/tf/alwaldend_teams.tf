resource "forgejo_team" "alwaldend_admins" {
  organization_id           = forgejo_organization.alwaldend.id
  name                      = "alwaldend_admins"
  permission                = "admin"
  can_create_org_repo       = true
  description               = "Admins"
  includes_all_repositories = true
  units_map = {
    "repo.code"       = "admin"
    "repo.actions"    = "admin"
    "repo.ext_issues" = "admin"
    "repo.ext_wiki"   = "admin"
    "repo.issues"     = "admin"
    "repo.packages"   = "admin"
    "repo.projects"   = "admin"
    "repo.pulls"      = "admin"
    "repo.releases"   = "admin"
    "repo.wiki"       = "admin"
  }
}

resource "forgejo_team" "alwaldend_devs" {
  organization_id           = forgejo_organization.alwaldend.id
  name                      = "alwaldend_devs"
  permission                = "write"
  description               = "Developers"
  can_create_org_repo       = false
  includes_all_repositories = false
  units_map = {
    "repo.code"     = "write"
    "repo.issues"   = "write"
    "repo.projects" = "write"
    "repo.pulls"    = "write"
    "repo.wiki"     = "write"
  }
}

locals {
  alwaldend_admin_entity_ids = toset(concat(
    tolist(data.vault_identity_group.forgejo_access["admins"].member_entity_ids),
    flatten([for group in data.vault_identity_group.forgejo_admin_members : tolist(group.member_entity_ids)]),
  ))
  alwaldend_admins = {
    for name, entity in local.vault_user_entities : name => entity
    if contains(local.alwaldend_admin_entity_ids, entity.entity_id)
  }
  alwaldend_package_writers = {
    for name, entity in local.vault_user_entities : name => entity
    if contains(data.vault_identity_group.forgejo_access["package_writers"].member_entity_ids, entity.entity_id)
  }
}

resource "forgejo_team_member" "alwaldend_admins" {
  for_each = local.alwaldend_admins
  team_id  = forgejo_team.alwaldend_admins.id
  user     = forgejo_user.vault[each.key].login
}

resource "forgejo_team" "alwaldend_package_writers" {
  organization_id           = forgejo_organization.alwaldend.id
  name                      = "alwaldend_package_writers"
  permission                = "write"
  can_create_org_repo       = false
  includes_all_repositories = false
  units_map = {
    "repo.packages" = "write"
  }
}

resource "forgejo_team_member" "alwaldend_package_writers" {
  for_each = local.alwaldend_package_writers
  team_id  = forgejo_team.alwaldend_package_writers.id
  user     = forgejo_user.vault[each.key].login
}
