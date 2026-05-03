module "src_infra_yandex_cloud_org1_tf_backend" {
    source = "../../../../projects/tf_modules/tf_backend"
    name = "src-infra-yandex-cloud-org1-tf-backend"
    bucket_prefix = "com-alwaldend"
    service_accounts = [var.service_account_id]
    folder_id = var.folder_id
    secret_name = "alwaldend.com/vault1/approles/src_infra_yandex_cloud_org1/bucket"
}
