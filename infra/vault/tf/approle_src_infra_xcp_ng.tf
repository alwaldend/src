module "src_infra_xcp_ng_approle" {
  source = "../../../projects/tf_modules/vault_approle"
  name   = "src_infra_xcp_ng"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  secrets                  = vault_mount.secrets.path
  backend                  = vault_auth_backend.approle.path
  backend_accessor         = vault_auth_backend.approle.accessor
  disable_yc_folder_policy = true
}

module "src_infra_xcp_ng_pki_server" {
  source  = "../../../projects/tf_modules/vault_pki_server"
  backend = module.pki_ica_servers.backend
  name    = "src_infra_xcp_ng_pki_server"
  allowed_domains = [
    "xoa.xcp-ng.alwaldend.com",
    "host1.xoa.xcp-ng.alwaldend.com",
    "host1.xcp-ng.alwaldend.com",
    "xcp-ng.alwaldend.com",
  ]
  eab_new_member_group_ids = [module.src_infra_xcp_ng_approle.group_id]
}

module "src_infra_xcp_ng_provider" {
  source      = "../../../projects/tf_modules/vault_oidc_provider"
  name        = "src_infra_xcp_ng_provider"
  issuer_host = var.vault_host
  scopes_supported = [
    vault_identity_oidc_scope.user.name,
    vault_identity_oidc_scope.groups.name,
  ]
  group_ids = [
    vault_identity_group.src_infra_xcp_ng_users.id,
  ]
  allowed_read_clients_group_ids = [
    module.src_infra_xcp_ng_approle.group_id,
  ]
  redirect_urls = [
    "${var.xoa_url}/signin/oidc/callback",
  ]
}

resource "vault_identity_group" "src_infra_xcp_ng_users" {
  name = "src_infra_xcp_ng_users"
  type = "internal"
  member_group_ids = [
    vault_identity_group.dev.id,
    vault_identity_group.approles.id,
    vault_identity_group.src_infra_xcp_ng_admins.id,
  ]
  metadata = {
    comment = "Users allowed to log in to Xen Orchestra with OIDC"
  }
}

resource "vault_identity_group" "src_infra_xcp_ng_admins" {
  name = "src_infra_xcp_ng_admins"
  type = "internal"
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  member_group_ids = [
    module.src_infra_xcp_ng_approle.group_id,
  ]
  metadata = {
    comment = "Xen Orchestra administrators"
  }
}
