terraform {
  required_providers {
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "0.203.0"
    }
  }
  backend "http" {
  }
}

provider "yandex" {
}
