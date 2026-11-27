terraform {
  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = "3.3.0"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
  }
}

provider "helm" {
  kubernetes = {
    insecure = false
  }
}

locals {
  namespace       = local.operator_issuer.namespace
  mirror_host     = "forgejo.alwaldend.com"
  mirror_owner    = "alwaldend"
  registry        = "${local.mirror_host}/${local.mirror_owner}"
  helm_repository = "https://${local.mirror_host}/api/packages/${local.mirror_owner}/helm"
  operator        = jsondecode(file("${path.module}/../../../third_party/io_controlplane_flux_operator/artifacts.json"))
  flux            = jsondecode(file("${path.module}/../../../third_party/io_fluxcd_flux/artifacts.json"))
  flux_cli        = one([for image in local.flux.images : image if endswith(image.repository, "/flux-cli")])
  flux_version    = trimprefix(local.flux_cli.tag, "v")
}

resource "helm_release" "operator" {
  name             = "flux-operator"
  namespace        = local.namespace
  create_namespace = true
  repository       = local.helm_repository
  chart            = one(local.operator.charts).name
  version          = one(local.operator.charts).version
  wait             = true
  atomic           = true

  values = [yamlencode({
    installCRDs = true
    image = {
      repository = "${local.registry}/flux-operator"
      tag        = "${one(local.operator.images).tag}@${one(local.operator.images).digest}"
    }
    extraEnvs = [{
      name  = "SSL_CERT_FILE"
      value = "/etc/flux-operator/vault/ca.crt"
    }]
    extraVolumeMounts = [{
      name      = "vault-ca"
      mountPath = "/etc/flux-operator/vault"
      readOnly  = true
    }]
    extraVolumes = [{
      name      = "vault-ca"
      configMap = { name = "flux-operator-vault-ca" }
    }]
    web = {
      enabled = true
      rbac    = { createRoles = false }
      config = {
        baseURL = "https://${local.operator_hostname}"
        authentication = {
          type = "OAuth2"
          oauth2 = {
            provider  = "OIDC"
            issuerURL = "${local.vault_url}/v1/identity/oidc/provider/src_infra_flux_provider"
            clientID  = var.operator_oidc_client_id
            scopes    = ["openid", "user", "groups", "email"]
            impersonation = {
              username = "'flux-ui:' + claims.sub"
              groups   = "claims.groups.map(group, 'flux-ui:' + group)"
            }
          }
        }
      }
      ingress = {
        enabled   = true
        className = "traefik"
        annotations = {
          "cert-manager.io/common-name"                      = local.operator_hostname
          "cert-manager.io/issuer"                           = local.operator_issuer.name
          "cert-manager.io/duration"                         = "168h"
          "traefik.ingress.kubernetes.io/router.entrypoints" = "websecure"
        }
        hosts = [{
          host  = local.operator_hostname
          paths = [{ path = "/", pathType = "Prefix" }]
        }]
        tls = [{
          hosts      = [local.operator_hostname]
          secretName = "flux-operator-tls"
        }]
      }
    }
  })]

  set_sensitive = [{
    name  = "web.config.authentication.oauth2.clientSecret"
    type  = "string"
    value = var.operator_oidc_client_secret
  }]

  depends_on = [helm_release.operator_ingress]
}

resource "helm_release" "instance" {
  name      = "flux-instance"
  namespace = local.namespace
  chart     = "${path.module}/flux-instance"
  version   = yamldecode(file("${path.module}/flux-instance/Chart.yaml")).version
  atomic    = true

  values = [yamlencode({
    distribution = {
      version  = local.flux_version
      registry = local.registry
      variant  = "upstream-alpine"
    }
  })]

  depends_on = [helm_release.operator]
}
