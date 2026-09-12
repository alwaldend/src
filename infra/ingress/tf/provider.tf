terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "5.22.0"
    }
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "0.201.0"
    }
    routeros = {
      source  = "terraform-routeros/routeros"
      version = "1.99.1"
    }
  }
  backend "http" {
  }
}

provider "yandex" {
}

provider "routeros" {
}
