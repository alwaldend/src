module "src_infra_dns_approle" {
    source = "../../../../projects/tf_modules/vault_approle"
    name = "src_infra_dns"
    member_entity_ids = [data.vault_identity_entity.simeonwarren.id]
    backend = vault_auth_backend.approle.path
    policy = <<EOT
        path "secrets/data/cloudflare.com/dns_token" {
            capabilities = ["read"]
        }
EOT
}
