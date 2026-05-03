module "src_infra_dc1_vault_tf_backend" {
    source = "../../../../projects/tf_modules/tf_backend"
    name = "src-infra-dc1-vault-tf-backend"
    bucket_prefix = "com-alwaldend"
    service_accounts = [var.service_account_id]
    folder_id = var.folder_id
    secret_name = "${module.src_infra_dc1_vault_approle.secret_path}/bucket"
}
