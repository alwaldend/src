terraform {
  required_providers {
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "0.203.0"
    }
    routeros = {
      source  = "terraform-routeros/routeros"
      version = "1.99.1"
    }
  }
  backend "http" {
  }
}

provider "yandex" {
}

provider "routeros" {
}
