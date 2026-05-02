module "src_infra_yandex_cloud_org1_folder" {
    source = "../../../../projects/tf_modules/yc_folder"
    name = "src-infra-yandex-cloud-org1"
    secret_name = "yandex.cloud/org1/folders/src-infra-yandex-cloud-org1"
    cloud_id = var.cloud_id
}

module "src_infra_dc1_vault_folder" {
    source = "../../../../projects/tf_modules/yc_folder"
    name = "src-infra-dc1-vault"
    secret_name = "yandex.cloud/org1/folders/src-infra-dc1-vault"
    cloud_id = var.cloud_id
}
