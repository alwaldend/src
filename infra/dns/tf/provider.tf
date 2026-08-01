terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "5.22.0"
    }
  }
  backend "http" {
  }
}

provider "cloudflare" {
}
