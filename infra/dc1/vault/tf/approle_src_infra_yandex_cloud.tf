module "src_infra_yandex_cloud_org1_approle" {
    source = "../../../../projects/tf_modules/vault_approle"
    name = "src_infra_yandex_cloud_org1"
    member_entity_ids = [data.vault_identity_entity.simeonwarren.id]
    secrets =  vault_mount.secrets.path
    policies = [vault_policy.tf_token.name]
    backend = vault_auth_backend.approle.path
    disable_yc_folder_policy = true
    policy = <<EOT
        path "${vault_mount.secrets.path}/data/yandex.cloud/org1/*" {
            capabilities = ["create", "update", "patch", "read", "delete"]
        }
        path "${vault_mount.secrets.path}/metadata/yandex.cloud/org1/*" {
            capabilities = ["list", "read", "update"]
        }
EOT
}
