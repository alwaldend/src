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
  backend "s3" {
    key = "main.tfstate"
  }
}

provider "vault" {
}

provider "consul" {
  address    = "https://host1.consul1.dc1.alwaldend.com:8501"
  datacenter = "dc1"
}
