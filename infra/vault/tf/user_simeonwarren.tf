resource "vault_identity_entity" "simeonwarren" {
  name = "simeonwarren"
  metadata = {
    username  = "simeonwarren"
    email     = "simeonwarren@alwaldend.com"
    sshpubkey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINe+6zNbs5HKrdZnHgeik/y7v5QYQr1cLYqfdCQ49GbM openpgp:0x485AEA03"
  }
}

resource "vault_identity_entity_alias" "userpass_simeonwarren" {
  name           = vault_identity_entity.simeonwarren.name
  mount_accessor = vault_auth_backend.userpass.accessor
  canonical_id   = vault_identity_entity.simeonwarren.id
}

module "users_simeonwarren_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "user_${vault_identity_entity.simeonwarren.name}"
  policies = [
    module.users_simeonwarren_ssh.policy,
  ]
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets          = vault_mount.secrets.path
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}

resource "vault_identity_entity_alias" "cert_simeonwarren" {
  name           = "simeonwarren.users.alwaldend.com"
  mount_accessor = vault_auth_backend.cert.accessor
  canonical_id   = vault_identity_entity.simeonwarren.id
}

resource "vault_cert_auth_backend_role" "cert_simeonwarren" {
  name                 = "simeonwarren"
  certificate          = file("${path.module}/../../../data/ssl/alwaldend.com/simeonwarren.crt")
  backend              = vault_auth_backend.cert.path
  allowed_common_names = ["simeonwarren.users.alwaldend.com"]
  token_ttl            = local.day_in_seconds
  token_max_ttl        = local.day_in_seconds
}

module "users_simeonwarren_ssh" {
  source          = "../../../projects/tf_modules/vault_ssh_server_role"
  backend         = vault_mount.ssh_servers.path
  name            = "users_simeonwarren_ssh"
  allowed_domains = "simeonwarren.users.alwaldend.com"
}

module "users_simeonwarren_pki_server" {
  source                   = "../../../projects/tf_modules/vault_pki_server"
  backend                  = module.pki_ica_servers.backend
  name                     = "users_simeonwarren_pki_server"
  allowed_domains          = ["simeonwarren.users.alwaldend.com"]
  eab_new_member_group_ids = [module.users_simeonwarren_approle.group_id]
  client_flag              = true
}

module "users_simeonwarren_provider" {
  source      = "../../../projects/tf_modules/vault_oidc_provider"
  name        = "users_simeonwarren_provider"
  issuer_host = var.vault_host
  scopes_supported = [
    vault_identity_oidc_scope.user.name,
    vault_identity_oidc_scope.groups.name,
    vault_identity_oidc_scope.ssh.name,
    vault_identity_oidc_scope.email.name,
  ]
  group_ids = [
    module.users_simeonwarren_approle.group_id,
  ]
  allowed_read_clients_group_ids = [
    module.users_simeonwarren_approle.group_id,
  ]
  redirect_urls = [
    "https://opencode.simeonwarren.users.alwaldend.com/traefik/oidc/callback",
  ]
}
