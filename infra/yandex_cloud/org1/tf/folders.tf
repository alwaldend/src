module "approle_folders" {
  for_each    = local.approles
  source      = "../../../../projects/tf_modules/yc_folder"
  name        = replace(each.value.name, "_", "-")
  secret_name = format("yandex.cloud/org1/folders/%s", replace(each.value.name, "_", "-"))
  cloud_id    = var.cloud_id
}
