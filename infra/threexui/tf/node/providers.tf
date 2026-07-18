terraform {
  required_providers {
    threexui = {
      source  = "batonogov/threexui"
      version = "3.20.0"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
  }
}
