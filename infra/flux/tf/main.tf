terraform {
  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = "3.3.0"
    }
  }
}

provider "helm" {
  kubernetes = {
    insecure = false
  }
}

locals {
  namespace    = "flux-system"
  registry     = "forgejo.alwaldend.com/alwaldend"
  operator     = jsondecode(file("${path.module}/../../../third_party/io_controlplane_flux_operator/artifacts.json"))
  flux         = jsondecode(file("${path.module}/../../../third_party/io_fluxcd_flux/artifacts.json"))
  flux_cli     = one([for image in local.flux.images : image if endswith(image.repository, "/flux-cli")])
  flux_version = trimprefix(local.flux_cli.tag, "v")
}

resource "helm_release" "operator" {
  name             = "flux-operator"
  namespace        = local.namespace
  create_namespace = true
  repository       = "https://forgejo.alwaldend.com/api/packages/alwaldend/helm"
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
    web = {
      enabled = false
    }
  })]
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
