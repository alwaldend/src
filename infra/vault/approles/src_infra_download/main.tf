terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
  }
  backend "http" {}
}

provider "vault" {}

data "vault_auth_backend" "approle" {
  path = "approle"
}

data "vault_identity_entity" "operator" {
  entity_name = "simeonwarren"
}

module "dns" {
  source  = "../../tf/dns_access"
  name    = "src_infra_download"
  secrets = "secrets"
  views   = ["dc1", "global"]
}

module "approle" {
  source                   = "../../../../projects/tf_modules/vault_approle"
  name                     = "src_infra_download"
  backend                  = data.vault_auth_backend.approle.path
  backend_accessor         = data.vault_auth_backend.approle.accessor
  secrets                  = "secrets"
  member_entity_ids        = [data.vault_identity_entity.operator.entity_id]
  disable_yc_folder_policy = false
  policies                 = [module.dns.policy_name, module.ssh.policy]
}

module "ssh" {
  source          = "../../../../projects/tf_modules/vault_ssh_server_role"
  backend         = "ssh/servers"
  name            = "src_infra_download_ssh"
  allowed_domains = "download.alwaldend.com"
}

module "pki_server" {
  source                   = "../../../../projects/tf_modules/vault_pki_server"
  backend                  = "pki/ica_servers"
  name                     = "src_infra_download_pki_server"
  allowed_domains          = ["alwaldend.com", "download.alwaldend.com", "www.alwaldend.com"]
  allow_subdomains         = false
  max_ttl                  = 7 * 24 * 60 * 60
  eab_new_member_group_ids = [module.approle.group_id]
}

output "entity_id" {
  description = "Component identity for infrastructure provider discovery"
  value       = module.approle.entity_id
}

output "group_id" {
  description = "Component group for Ansible client signing"
  value       = module.approle.group_id
}
