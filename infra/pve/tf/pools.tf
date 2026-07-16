locals {
  entities_by_name = { for _, entity in data.vault_identity_entity.approles_entities : entity.entity_name => entity }
}

data "vault_identity_group" "approles" {
  group_name = "approles"
}

data "vault_identity_entity" "approles_entities" {
  for_each  = { for id in data.vault_identity_group.approles.member_entity_ids : id => "" }
  entity_id = each.key
}

resource "proxmox_pool" "approles" {
  for_each = local.entities_by_name
  poolid   = each.key
  comment  = "Pool for ${each.key} (${each.value.entity_id}) approle, managed by terraform (//infra/dc1/pve1)"
}

resource "proxmox_pool" "templates" {
  poolid  = "templates"
  comment = "Template pool. Managed by terraform (//infra/dc1/pve1)"
}
