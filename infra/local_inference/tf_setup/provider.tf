terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "5.22.0"
    }
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "0.203.0"
    }
  }
  backend "http" {}
}

provider "yandex" {
  folder_id = var.folder_id
  zone      = var.zone
}

provider "cloudflare" {}
