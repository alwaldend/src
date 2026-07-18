locals {
  users_by_name = {
    for _, entity in data.vault_identity_entity.threexui_users : entity.entity_name => entity
  }
  admins_by_name = {
    for _, entity in data.vault_identity_entity.threexui_admins : entity.entity_name => entity
  }
}

data "vault_identity_group" "threexui_users" {
  group_name = "src_infra_threexui_users"
}

data "vault_identity_entity" "threexui_users" {
  for_each = {
    for id in data.vault_identity_group.threexui_users.member_entity_ids : id => ""
  }
  entity_id = each.key
}

data "vault_identity_group" "threexui_admins" {
  group_name = "src_infra_threexui_admins"
}

data "vault_identity_entity" "threexui_admins" {
  for_each = {
    for id in data.vault_identity_group.threexui_admins.member_entity_ids : id => ""
  }
  entity_id = each.key
}
