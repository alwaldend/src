terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
    harbor = {
      source  = "goharbor/harbor"
      version = "3.12.0"
    }
  }
  backend "http" {
  }
}

provider "vault" {
}

provider "harbor" {
}
