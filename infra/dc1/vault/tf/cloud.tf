module "src_infra_dc1_vault_tf_backend" {
    source = "../../../../projects/tf_modules/tf_backend"
    name = "src-infra-dc1-vault-tf-backend"
    bucket_prefix = "com-alwaldend"
    service_accounts = [var.service_account_id]
    folder_id = var.folder_id
    secret_name = "${module.src_infra_dc1_vault_approle.secret_path}/bucket"
}

module "backup_bucket" {
    source = "../../../../projects/tf_modules/backup_bucket"
    name = "src-infra-dc1-vault-backup"
    bucket_prefix = "com-alwaldend"
    service_accounts = [var.service_account_id]
    folder_id = var.folder_id
    secret_name = "${module.src_infra_dc1_vault_approle.secret_path}/backup_bucket"
}
