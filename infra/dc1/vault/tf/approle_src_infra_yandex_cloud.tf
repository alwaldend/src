module "src_infra_yandex_cloud_org1_approle" {
    source = "../../../../projects/tf_modules/vault_approle"
    name = "src_infra_yandex_cloud_org1"
    member_entity_ids = [data.vault_identity_entity.simeonwarren.id]
    backend = vault_auth_backend.approle.path
    policy = <<EOT
        path "secrets/data/yandex.cloud/org1/*" {
            capabilities = ["create", "update", "patch", "read", "delete"]
        }
        path "secrets/metadata/yandex.cloud/org1/*" {
            capabilities = ["list"]
        }
EOT
}
