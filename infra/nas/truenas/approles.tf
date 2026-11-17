data "vault_identity_group" "truenas" {
  for_each   = toset(["truenas_admins", "truenas_user"])
  group_name = each.key
}

data "vault_identity_entity" "approles" {
  for_each  = toset(flatten([for group in data.vault_identity_group.truenas : group.member_entity_ids]))
  entity_id = each.key
}

locals {
  nas_password_version = 2
  nas_memberships = {
    for id, entity in data.vault_identity_entity.approles : entity.entity_name => [
      for name, group in data.vault_identity_group.truenas : name
      if contains(group.member_entity_ids, id)
    ]
  }
  nas_primary_groups = {
    for name, groups in local.nas_memberships : name => contains(groups, "truenas_admins") ? "truenas_admins" : "truenas_user"
  }
  approles = toset(keys(local.nas_memberships))
  nas_usernames = {
    for name in local.approles : name => length(name) <= 32 ? name : "${substr(name, 0, 23)}_${substr(sha256(name), 0, 8)}"
  }
}

resource "random_password" "approles" {
  for_each = local.approles
  length   = 64
  special  = false
  keepers  = { version = local.nas_password_version }
}

resource "truenas_group" "access" {
  for_each = data.vault_identity_group.truenas
  name     = each.key
  smb      = true
}

resource "truenas_privilege" "administrators" {
  name         = "truenas_admins"
  local_groups = [truenas_group.access["truenas_admins"].gid]
  roles        = ["FULL_ADMIN"]
  web_shell    = false
}

resource "truenas_user" "approles" {
  for_each            = local.approles
  username            = local.nas_usernames[each.key]
  full_name           = each.key
  password            = random_password.approles[each.key].result
  password_wo_version = local.nas_password_version
  group_create        = false
  group               = truenas_group.access[local.nas_primary_groups[each.key]].id
  groups = [
    for group in local.nas_memberships[each.key] : truenas_group.access[group].id
    if group != local.nas_primary_groups[each.key]
  ]
  home                   = "/var/empty"
  shell                  = "/usr/sbin/nologin"
  smb                    = true
  ssh_password_enabled   = false
  sudo_commands          = []
  sudo_commands_nopasswd = []

  lifecycle {
    precondition {
      condition     = length(distinct(values(local.nas_usernames))) == length(local.nas_usernames)
      error_message = "NAS account names must be unique after shortening long AppRole names."
    }
  }
}

resource "vault_kv_secret_v2" "nas_accounts" {
  for_each = local.approles
  mount    = "secrets"
  name     = "alwaldend.com/vault1/shared/read/approle/${each.key}/nas"
  data_json = jsonencode({
    username = truenas_user.approles[each.key].username
    password = random_password.approles[each.key].result
  })
}
