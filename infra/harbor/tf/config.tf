resource "harbor_config_auth" "oidc" {
  auth_mode          = "oidc_auth"
  primary_auth_mode  = false
  oidc_name          = "vault"
  oidc_endpoint      = "https://vault.dc1.alwaldend.com:8200/v1/identity/oidc/provider/src_infra_harbor_provider"
  oidc_client_id     = data.vault_identity_oidc_client_creds.oidc.client_id
  oidc_client_secret = data.vault_identity_oidc_client_creds.oidc.client_secret
  oidc_group_filter  = ""
  oidc_groups_claim  = "groups"
  oidc_admin_group   = "src_infra_harbor_admins"
  oidc_scope         = "openid,offline_access,user,email,groups"
  oidc_verify_cert   = false // Ca issues
  oidc_auto_onboard  = true
  oidc_logout        = true
  oidc_user_claim    = "username"
}

data "vault_identity_oidc_client_creds" "oidc" {
  name = "src_infra_harbor_provider"
}
