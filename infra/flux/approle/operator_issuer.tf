locals {
  operator_hostname = jsondecode(file("${path.module}/../dns/dnsconfig.json")).records.operator.CNAME.name
  operator_issuer   = jsondecode(file("${path.module}/../operator-issuer.json"))
}

resource "vault_auth_backend" "operator_issuer" {
  type = "kubernetes"
  path = local.operator_issuer.vaultAuth
}

resource "vault_policy" "operator_issuer_config" {
  name   = "${local.operator_issuer.vaultRole}_issuer_config"
  policy = <<-EOT
    path "auth/${vault_auth_backend.operator_issuer.path}/config" {
      capabilities = ["create", "read", "update"]
    }
  EOT
}

resource "vault_pki_secret_backend_role" "operator_issuer" {
  backend                     = local.operator_issuer.vaultPKI
  name                        = local.operator_issuer.vaultRole
  allowed_domains             = [local.operator_hostname]
  allow_bare_domains          = true
  allow_subdomains            = false
  allow_wildcard_certificates = false
  allow_localhost             = false
  allow_ip_sans               = false
  client_flag                 = false
  key_type                    = "rsa"
  key_bits                    = 2048
  ttl                         = 604800
  max_ttl                     = 604800
}

resource "vault_policy" "operator_sign" {
  name   = "${local.operator_issuer.vaultRole}_sign"
  policy = <<-EOT
    path "${vault_pki_secret_backend_role.operator_issuer.backend}/issuer/default/sign/${vault_pki_secret_backend_role.operator_issuer.name}" {
      capabilities = ["update"]
    }
  EOT
}

resource "vault_kubernetes_auth_backend_role" "operator_issuer" {
  backend                          = vault_auth_backend.operator_issuer.path
  role_name                        = local.operator_issuer.vaultRole
  bound_service_account_names      = [local.operator_issuer.serviceAccount]
  bound_service_account_namespaces = [local.operator_issuer.namespace]
  audience                         = "vault://${local.operator_issuer.namespace}/${local.operator_issuer.name}"
  token_policies                   = [vault_policy.operator_sign.name]
  token_no_default_policy          = true
  token_ttl                        = 300
  token_max_ttl                    = 300
}
