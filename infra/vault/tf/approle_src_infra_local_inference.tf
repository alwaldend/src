module "src_infra_local_inference_approle" {
  source   = "../../../projects/tf_modules/vault_approle"
  name     = "src_infra_local_inference"
  policies = [module.dns_access["src_infra_local_inference"].policy_name]
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets          = vault_mount.secrets.path
  backend          = vault_auth_backend.approle.path
  backend_accessor = vault_auth_backend.approle.accessor
}

module "src_infra_local_inference_pki_server" {
  source  = "../../../projects/tf_modules/vault_pki_server"
  name    = "src_infra_local_inference_pki_server"
  backend = module.pki_ica_servers.backend
  allowed_domains = [
    jsondecode(file("${path.module}/../../local_inference/dnsconfig.json")).records.t3code.CNAME.name,
  ]
  allow_subdomains         = false
  client_flag              = false
  eab_new_member_group_ids = [module.src_infra_local_inference_approle.group_id]
}
