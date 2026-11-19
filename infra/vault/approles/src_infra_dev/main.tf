data "vault_auth_backend" "approle" {
  path = "approle"
}

data "vault_identity_entity" "operator" {
  entity_name = "simeonwarren"
}

module "dns" {
  source  = "../../tf/dns_access"
  name    = "src_infra_dev"
  secrets = "secrets"
  views   = ["dc1", "global"]
}

module "approle" {
  source                   = "../../../../projects/tf_modules/vault_approle"
  name                     = "src_infra_dev"
  disable_yc_folder_policy = true
  member_entity_ids        = [data.vault_identity_entity.operator.entity_id]
  secrets                  = "secrets"
  backend                  = data.vault_auth_backend.approle.path
  backend_accessor         = data.vault_auth_backend.approle.accessor
  policies = [
    module.ssh.policy,
    module.dns.policy_name,
    vault_policy.runner.name,
  ]
}

resource "vault_policy" "runner" {
  name   = "src_infra_dev_runner"
  policy = <<-EOT
    path "secrets/data/alwaldend.com/vault1/approles/src_infra_dc1_forgejo1/runners/secure1" {
      capabilities = ["read"]
    }
    path "secrets/data/alwaldend.com/vault1/approles/src_infra_threexui/http_proxies" {
      capabilities = ["read"]
    }
  EOT
}

locals {
  dev_dns = jsondecode(file("${path.module}/../../../dev/dns/dnsconfig.json"))
}

module "ssh" {
  source           = "../../../../projects/tf_modules/vault_ssh_server_role"
  backend          = "ssh/servers"
  name             = "src_infra_dev_ssh"
  allow_subdomains = false
  allowed_domains = join(",", [
    for key in ["secure1", "dev1_dc1"] : "${local.dev_dns.records[key].A.name}.alwaldend.com"
  ])
}

module "pki_server" {
  source                   = "../../../../projects/tf_modules/vault_pki_server"
  backend                  = "pki/ica_servers"
  name                     = "src_infra_dev_pki_server"
  allow_subdomains         = false
  eab_new_member_group_ids = [module.approle.group_id]
  allowed_domains = [
    for key in ["t3code_dc1", "comfyui_dc1"] : "${local.dev_dns.records[key].A.name}.alwaldend.com"
  ]
}
