terraform {
  required_providers {
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
