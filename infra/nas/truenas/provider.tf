terraform {
  required_version = ">= 1.11.0"
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "3.9.1"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
    truenas = {
      source  = "truenas/truenas"
      version = "1.9.3"
    }
  }
}

provider "vault" {
  skip_child_token = true
}

# Connection and API key arrive through the provider's native environment.
provider "truenas" {
  insecure = false
}
