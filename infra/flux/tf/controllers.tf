locals {
  cert_manager = jsondecode(file("${path.module}/../../../third_party/io_cert_manager_cert_manager/artifacts.json"))
  traefik      = jsondecode(file("${path.module}/../../../third_party/io_traefik_traefik/artifacts.json"))
  cert_manager_images = {
    for image in local.cert_manager.images : basename(image.repository) => {
      repository = "${local.registry}/${basename(image.repository)}"
      tag        = image.tag
      digest     = image.digest
    }
  }
}

resource "helm_release" "cert_manager" {
  name             = "cert-manager"
  namespace        = "cert-manager"
  create_namespace = true
  repository       = local.helm_repository
  chart            = one(local.cert_manager.charts).name
  version          = one(local.cert_manager.charts).version
  atomic           = true
  wait             = true

  values = [yamlencode({
    crds = {
      enabled = true
    }
    image = local.cert_manager_images["cert-manager-controller"]
    webhook = {
      image = local.cert_manager_images["cert-manager-webhook"]
    }
    cainjector = {
      image = local.cert_manager_images["cert-manager-cainjector"]
    }
    acmesolver = {
      image = local.cert_manager_images["cert-manager-acmesolver"]
    }
    startupapicheck = {
      image = local.cert_manager_images["cert-manager-startupapicheck"]
    }
  })]
}

resource "helm_release" "traefik" {
  name             = "traefik"
  namespace        = "traefik"
  create_namespace = true
  repository       = local.helm_repository
  chart            = one(local.traefik.charts).name
  version          = one(local.traefik.charts).version
  atomic           = true
  wait             = true

  values = [yamlencode({
    image = {
      registry   = local.mirror_host
      repository = "${local.mirror_owner}/${basename(one(local.traefik.images).repository)}"
      tag        = "${one(local.traefik.images).tag}@${one(local.traefik.images).digest}"
    }
    service = {
      spec = {
        type = "ClusterIP"
      }
    }
  })]
}
