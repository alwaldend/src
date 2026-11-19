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
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
    proxmox = {
      source  = "Telmate/proxmox"
      version = "3.0.2-rc07"
    }
  }
  backend "http" {
  }
}

provider "vault" {
}

provider "proxmox" {
}
