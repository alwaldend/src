terraform {
  required_providers {
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
