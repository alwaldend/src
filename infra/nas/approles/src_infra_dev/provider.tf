terraform {
  required_providers {
    truenas = {
      source  = "truenas/truenas"
      version = "1.9.3"
    }
  }
}

provider "truenas" {
  insecure = false
}
