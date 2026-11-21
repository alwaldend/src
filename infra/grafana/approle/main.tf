locals {
  hostname = jsondecode(file("${path.module}/../dns/dnsconfig.json")).records.local.A.name
  name     = "src_infra_grafana"
}

data "vault_auth_backend" "approle" {
  path = "approle"
}

data "vault_identity_group" "operators" {
  group_name = "global_admins"
}

module "approle" {
  source                   = "../../../projects/tf_modules/vault_approle"
  name                     = local.name
  member_group_ids         = [data.vault_identity_group.operators.group_id]
  policies                 = [module.dns.policy_name, module.ssh.policy]
  secrets                  = "secrets"
  disable_yc_folder_policy = true
  backend                  = data.vault_auth_backend.approle.path
  backend_accessor         = data.vault_auth_backend.approle.accessor
}

module "dns" {
  source  = "../../vault/tf/dns_access"
  name    = local.name
  secrets = "secrets"
  views   = ["dc1", "global"]
}

module "pki_server" {
  source                   = "../../../projects/tf_modules/vault_pki_server"
  backend                  = "pki/ica_servers"
  name                     = "${local.name}_pki_server"
  allowed_domains          = [local.hostname]
  allow_subdomains         = false
  eab_new_member_group_ids = [module.approle.group_id]
}

module "ssh" {
  source           = "../../../projects/tf_modules/vault_ssh_server_role"
  backend          = "ssh/servers"
  name             = "${local.name}_ssh"
  allowed_domains  = join(",", [split(".", local.hostname)[0], local.hostname])
  allow_subdomains = false
}

module "oidc" {
  source                         = "../../../projects/tf_modules/vault_oidc_provider"
  name                           = "${local.name}_provider"
  issuer_host                    = "vault.alwaldend.com:8200"
  scopes_supported               = ["user", "email", "groups"]
  group_ids                      = [data.vault_identity_group.operators.group_id]
  allowed_read_clients_group_ids = [module.approle.group_id]
  redirect_urls                  = ["https://${local.hostname}/login/generic_oauth"]
}
