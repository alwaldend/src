terraform {
  required_providers {
    external = {
      source  = "hashicorp/external"
      version = "2.3.5"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
    xenorchestra = {
      source  = "vatesfr/xenorchestra"
      version = "0.41.0"
    }
  }
  backend "http" {
  }
}

provider "xenorchestra" {
  url      = var.xoa_url
  insecure = var.xoa_insecure
}

provider "vault" {
}

variable "xoa_url" {
  description = "Xen Orchestra WebSocket endpoint."
  type        = string
  default     = "wss://xoa.xcp-ng.alwaldend.com"
}

variable "xoa_insecure" {
  description = "Explicit bootstrap override for an appliance without a valid certificate."
  type        = bool
  default     = false
}
