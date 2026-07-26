terraform {
  required_providers {
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "0.201.0"
    }
  }
  backend "http" {
  }
}

provider "yandex" {
}
