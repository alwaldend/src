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
