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

module "approle_folders" {
  for_each    = local.entities_by_name
  source      = "../../../../projects/tf_modules/yc_folder"
  name        = replace(each.key, "_", "-")
  secret_name = format("yandex.cloud/org1/folders/%s", replace(each.key, "_", "-"))
  cloud_id    = var.cloud_id
}
