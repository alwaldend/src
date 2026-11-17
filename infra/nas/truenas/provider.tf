terraform {
  required_version = ">= 1.11.0"
  required_providers {
    truenas = {
      source  = "truenas/truenas"
      version = "1.9.3"
    }
  }
}

# Connection and API key arrive through the provider's native environment.
provider "truenas" {
  insecure = false
}
