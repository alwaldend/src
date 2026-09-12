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
    proxmox = {
      source  = "Telmate/proxmox"
      version = "3.0.2-rc07"
    }
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "0.201.0"
    }
  }
  backend "http" {
  }
}

provider "yandex" {
}

provider "proxmox" {
  pm_minimum_permission_check = false
}
