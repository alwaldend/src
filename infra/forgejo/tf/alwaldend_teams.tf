resource "forgejo_team" "alwaldend_admins" {
  organization_id = forgejo_organization.alwaldend.id
  name            = "alwaldend_admins"
  permission      = "admin"
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

locals {
  alwaldend_admins = ["simeonwarren", "src_infra_dc1_forgejo1"]
  alwaldend_package_writers = ["src_third_party"]
}

resource "forgejo_team_member" "alwaldend_admins" {
  for_each = { for admin in local.alwaldend_admins : admin => "" }
  team_id  = forgejo_team.alwaldend_admins.id
  user     = each.key
}

resource "forgejo_team" "alwaldend_package_writers" {
  organization_id = forgejo_organization.alwaldend.id
  name            = "alwaldend_package_writers"
  permission = "write"
  units_map = {
    "repo.packages"   = "write"
  }
}

resource "forgejo_team_member" "alwaldend_package_writers" {
  for_each = { for login in local.alwaldend_package_writers : login => "" }
  team_id  = forgejo_team.alwaldend_package_writers.id
  user     = each.key
}
