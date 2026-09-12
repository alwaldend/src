terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "5.22.0"
    }
    proxmox = {
      source  = "Telmate/proxmox"
      version = "3.0.2-rc07"
    }
  }
  backend "http" {
  }
}

provider "proxmox" {
  pm_minimum_permission_check = false
}

provider "cloudflare" {}
