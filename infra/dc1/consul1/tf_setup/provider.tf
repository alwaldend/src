terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
    proxmox = {
      source  = "Telmate/proxmox"
      version = "3.0.2-rc07"
    }
  }
  backend "s3" {
    key = "setup.tfstate"
  }
}

provider "vault" {
}

provider "proxmox" {
  pm_minimum_permission_check = false
}
