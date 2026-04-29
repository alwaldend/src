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
}

provider "vault" {
}

provider "yandex" {
  token     = "auth_token_here"
  cloud_id  = "cloud_id_here"
  folder_id = "folder_id_here"
  zone      = "ru-central1-d"
}
