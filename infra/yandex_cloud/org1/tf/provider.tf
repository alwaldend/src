terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "0.203.0"
    }
  }
  backend "http" {
  }
}

provider "vault" {
}

provider "yandex" {
}
