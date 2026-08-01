terraform {
  required_providers {
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
