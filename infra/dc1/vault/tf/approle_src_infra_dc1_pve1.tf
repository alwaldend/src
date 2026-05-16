module "src_infra_dc1_pve1_approle" {
  source = "../../../../projects/tf_modules/vault_approle"
  name   = "src_infra_dc1_pve1"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets = vault_mount.secrets.path
  policies = [
    vault_policy.tf_token.name,
    vault_policy.ssh_servers_sign_dc1_pve1.name,
  ]
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}

resource "vault_ssh_secret_backend_role" "servers_dc1_pve1" {
  backend                 = vault_mount.ssh_servers.path
  name                    = "dc1_pve1"
  max_ttl                 = local.year_in_seconds
  ttl                     = local.year_in_seconds
  key_type                = "ca"
  allow_host_certificates = true
  allow_subdomains        = true
  allow_bare_domains      = true
  allowed_domains         = "pve1.dc1.alwaldend.com"
}

resource "vault_policy" "ssh_servers_sign_dc1_pve1" {
  name   = "ssh_servers_sign_dc1_pve1"
  policy = <<EOT
    path
    "${vault_mount.ssh_servers.path}/sign/${vault_ssh_secret_backend_role.servers_dc1_pve1.name}" {
      capabilities = ["update"]
    }
EOT
}
