terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "5.22.0"
    }
    routeros = {
      source  = "terraform-routeros/routeros"
      version = "1.99.1"
    }
    xenorchestra = {
      source  = "vatesfr/xenorchestra"
      version = "0.41.0"
    }
  }
  backend "http" {
  }
}

provider "xenorchestra" {
  url      = var.xoa_url
  insecure = var.xoa_insecure
}
