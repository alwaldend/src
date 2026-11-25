terraform {
  required_providers {
    threexui = {
      source  = "batonogov/threexui"
      version = "3.34.0"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
  }
  backend "http" {
  }
}

provider "threexui" {
  endpoint = var.xui_url
  username = var.xui_username
  password = var.xui_password
}

provider "vault" {
}
