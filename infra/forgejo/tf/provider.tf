terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
    forgejo = {
      source  = "svalabs/forgejo"
      version = "1.5.0"
    }
  }
  backend "http" {
  }
}

provider "vault" {
}

provider "forgejo" {
}
