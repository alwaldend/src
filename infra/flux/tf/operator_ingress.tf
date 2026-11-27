variable "kubernetes_ca" {
  type        = string
  description = "Public K3s API CA published by the bootstrap role."
}

variable "operator_oidc_client_id" {
  type        = string
  description = "Vault OIDC client ID for the Operator UI."
}

variable "operator_oidc_client_secret" {
  type        = string
  description = "Vault OIDC client secret for the Operator UI."
  sensitive   = true
}

provider "vault" {
  skip_child_token = true
}

locals {
  dns_records       = jsondecode(file("${path.module}/../dns/dnsconfig.json")).records
  operator_hostname = local.dns_records.operator.CNAME.name
  operator_issuer   = jsondecode(file("${path.module}/../operator-issuer.json"))
  ingress_access    = jsondecode(file("${path.module}/../ingress-access.json"))
  ingress_hosts     = yamldecode(file("${path.module}/../../ingress/ansible/inventory.yaml")).ingress.hosts
  vault_url         = "https://vault.alwaldend.com:8200"
  vault_ca_path     = "${path.module}/../../../data/ssl/alwaldend.com/pki_ca_servers.crt"
}

resource "vault_kubernetes_auth_backend_config" "operator_issuer" {
  backend              = local.operator_issuer.vaultAuth
  kubernetes_host      = "https://${local.dns_records.host1.A.address}:6443"
  kubernetes_ca_cert   = var.kubernetes_ca
  disable_local_ca_jwt = true
}

resource "helm_release" "operator_ingress" {
  name             = "flux-operator-ingress"
  namespace        = local.namespace
  create_namespace = true
  chart            = "${path.module}/operator-ingress"
  version          = yamldecode(file("${path.module}/operator-ingress/Chart.yaml")).version
  atomic           = true

  values = [
    jsonencode(local.operator_issuer),
    yamlencode({
      vaultURL             = local.vault_url
      vaultCA              = file(local.vault_ca_path)
      ingressSources       = concat([local.ingress_access.managementNetwork, local.ingress_access.workloadNetwork, "127.0.0.1/32", "::1/128"], [for host in values(local.ingress_hosts) : "${split("/", host.wg_ip)[0]}/32"])
      certManagerName      = helm_release.cert_manager.name
      certManagerNamespace = helm_release.cert_manager.namespace
    }),
  ]

  depends_on = [
    helm_release.traefik,
    vault_kubernetes_auth_backend_config.operator_issuer,
  ]
}
