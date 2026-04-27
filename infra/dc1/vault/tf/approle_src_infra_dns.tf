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

moved {
    from = vault_approle_auth_backend_role.src_infra_dns
    to = module.src_infra_dns_approle.vault_approle_auth_backend_role.role
}

moved {
    from = vault_policy.src_infra_dns
    to = module.src_infra_dns_approle.vault_policy.policy
}

moved {
    from = vault_policy.src_infra_dns_approle
    to = module.src_infra_dns_approle.vault_policy.policy_approle
}

moved {
    from = vault_identity_group.src_infra_dns
    to = module.src_infra_dns_approle.vault_identity_group.group
}
