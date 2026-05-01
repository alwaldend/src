module "src_infra_yandex_cloud_org1_tf_backend" {
    source = "../../../../projects/tf_modules/tf_backend"
    name = "src-infra-yandex-cloud-org1-tf-backend"
    bucket_prefix = "com-alwaldend"
    folder_id = var.folder_id
    secret_name = "yandex.cloud/org1/tf_backend/bucket"
}

module "src_infra_dc1_vault_tf_backend" {
    source = "../../../../projects/tf_modules/tf_backend"
    name = "src-infra-dc1-vault-tf-backend"
    bucket_prefix = "com-alwaldend"
    folder_id = var.folder_id
    secret_name = "yandex.cloud/org1/shared/read/approle/src_infra_dc1_vault/bucket"
}
