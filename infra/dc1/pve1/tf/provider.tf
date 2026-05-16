terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
  }
  backend "s3" {
  }
}

provider "vault" {
}
