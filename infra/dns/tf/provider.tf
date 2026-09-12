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
  }
  backend "http" {
  }
}

provider "cloudflare" {
}

provider "routeros" {
  alias    = "dns"
  hosturl  = var.dns_routeros_hosturl
  username = var.dns_routeros_username
  password = var.dns_routeros_password
}
