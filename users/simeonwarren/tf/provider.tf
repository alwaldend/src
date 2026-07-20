terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "0.201.0"
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

provider "yandex" {
}

provider "proxmox" {
  pm_minimum_permission_check = false
}
