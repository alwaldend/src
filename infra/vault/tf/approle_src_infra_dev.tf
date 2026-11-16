module "src_infra_dev_approle" {
  source                   = "../../../projects/tf_modules/vault_approle"
  name                     = "src_infra_dev"
  disable_yc_folder_policy = true
  member_entity_ids        = [vault_identity_entity.simeonwarren.id]
  secrets                  = vault_mount.secrets.path
  backend                  = vault_auth_backend.approle.path
  backend_accessor         = vault_auth_backend.approle.accessor
  policies = [
    module.src_infra_dev_ssh.policy,
    module.dns_access["src_infra_dev"].policy_name,
    vault_policy.src_infra_dev_runner.name,
  ]
}

resource "vault_policy" "src_infra_dev_runner" {
  name   = "src_infra_dev_runner"
  policy = <<-EOT
    path "${vault_mount.secrets.path}/data/alwaldend.com/vault1/approles/src_infra_dc1_forgejo1/runners/secure1" {
      capabilities = ["read"]
    }
  EOT
}

locals {
  dev_dns = jsondecode(file("${path.module}/../../dev/dnsconfig.json"))
}

module "src_infra_dev_ssh" {
  source           = "../../../projects/tf_modules/vault_ssh_server_role"
  backend          = vault_mount.ssh_servers.path
  name             = "src_infra_dev_ssh"
  allow_subdomains = false
  allowed_domains = join(",", [
    for key in ["secure1", "dev1_dc1"] : "${local.dev_dns.records[key].A.name}.alwaldend.com"
  ])
}

module "src_infra_dev_pki_server" {
  source                   = "../../../projects/tf_modules/vault_pki_server"
  backend                  = module.pki_ica_servers.backend
  name                     = "src_infra_dev_pki_server"
  allow_subdomains         = false
  eab_new_member_group_ids = [module.src_infra_dev_approle.group_id]
  allowed_domains = [
    for key in ["t3code_dc1", "comfyui_dc1"] : "${local.dev_dns.records[key].A.name}.alwaldend.com"
  ]
}
