data "vault_generic_secret" "oidc_client" {
  path = "identity/oidc/client/src_infra_xcp_ng_provider"
}

data "xenorchestra_pools" "all" {
}

locals {
  xo_config_binary = abspath("${path.module}/../cmd/xo_config/xo_config_/xo_config")
  oidc_authorization = {
    admin_group = "src_infra_xcp_ng_admins"
    user_group  = "src_infra_xcp_ng_users"
  }
  oidc_configuration = {
    discoveryURL = "https://vault.alwaldend.com:8200/v1/identity/oidc/provider/src_infra_xcp_ng_provider/.well-known/openid-configuration"
    clientID     = data.vault_generic_secret.oidc_client.data.client_id
    clientSecret = data.vault_generic_secret.oidc_client.data.client_secret
    advanced = {
      callbackURL   = "${trimsuffix(replace(var.xoa_url, "wss://", "https://"), "/")}/signin/oidc/callback"
      usernameField = "username"
      scope         = "user groups"
    }
  }
}

# The upstream provider has no plugin resource. Reconcile through XO's
# JSON-RPC API when configuration changes, without invoking a shell.
resource "terraform_data" "oidc" {
  triggers_replace = [
    sha256(jsonencode(local.oidc_configuration)),
    sha256(jsonencode(local.oidc_authorization)),
    filesha256(local.xo_config_binary),
  ]

  provisioner "local-exec" {
    interpreter = [local.xo_config_binary]
    command     = "apply-oidc"
    environment = {
      XOA_URL               = var.xoa_url
      XOA_INSECURE          = tostring(var.xoa_insecure)
      XO_OIDC_CONFIGURATION = jsonencode(local.oidc_configuration)
      XO_OIDC_AUTHORIZATION = jsonencode(local.oidc_authorization)
    }
    quiet = true
  }
}
