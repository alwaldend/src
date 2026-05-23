terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "2.9.0"
    }
  }
  backend "s3" {
    key = "main.tfstate"
  }
}

provider "vault" {
}
