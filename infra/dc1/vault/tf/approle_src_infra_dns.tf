resource "vault_identity_group" "src_infra_dns" {
  name     = "src_infra_dns"
  type     = "internal"
  policies = [vault_policy.src_infra_dns_approle.name]
  member_entity_ids = [vault_identity_entity.simeonwarren.id]

  metadata = {
    version = "2"
  }
}

resource "vault_approle_auth_backend_role" "src_infra_dns" {
  backend        = vault_auth_backend.approle.path
  role_name      = "src_infra_dns"
  role_id        = "src_infra_dns"
  secret_id_num_uses = 1
  secret_id_ttl = 3600
  token_policies = [vault_policy.src_infra_dns.name]
}

resource "vault_policy" "src_infra_dns_approle" {
  name = "src_infra_dns_approle"
  policy = <<EOT
    path "auth/${vault_auth_backend.approle.path}/role/${vault_approle_auth_backend_role.src_infra_dns.role_name}/secret-id" {
        capabilities = ["create"]
    }
EOT
}
resource "vault_policy" "src_infra_dns" {
  name = "src_infra_dns"
  policy = <<EOT
    path "secrets/data/cloudflare.com/dns_token" {
        capabilities = ["read"]
    }
EOT
}
