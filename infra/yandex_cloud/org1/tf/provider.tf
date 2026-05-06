terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "0.201.0"
    }
  }
  backend "s3" {
  }
}

provider "vault" {
}

provider "yandex" {
}
