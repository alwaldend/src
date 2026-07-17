terraform {
  required_providers {
    threexui = {
      source  = "batonogov/threexui"
      version = "3.20.0"
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
