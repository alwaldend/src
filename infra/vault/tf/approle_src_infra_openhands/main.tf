module "approle" {
  source                   = "../../../../projects/tf_modules/vault_approle"
  name                     = var.name
  member_entity_ids        = var.member_entity_ids
  secrets                  = var.secrets
  disable_yc_folder_policy = true
  policies = concat([
    module.ssh.policy,
  ], var.policies)
  backend          = var.approle_backend
  backend_accessor = var.approle_backend_accessor
}

module "ssh" {
  source          = "../../../../projects/tf_modules/vault_ssh_server_role"
  backend         = var.ssh_backend
  name            = "${var.name}_ssh"
  allowed_domains = var.ssh_allowed_domains
}

module "pki_server" {
  source                   = "../../../../projects/tf_modules/vault_pki_server"
  backend                  = var.pki_backend
  name                     = "${var.name}_pki_server"
  allowed_domains          = var.pki_allowed_domains
  eab_new_member_group_ids = [module.approle.group_id]
  client_flag              = true
}
