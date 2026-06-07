terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
    consul = {
      source  = "hashicorp/consul"
      version = "2.23.0"
    }
  }
  backend "http" {
  }
}

provider "vault" {
}

provider "consul" {
  datacenter = "dc1"
}
