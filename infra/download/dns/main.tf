terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "5.22.0"
    }
    routeros = {
      source  = "terraform-routeros/routeros"
      version = "1.99.1"
    }
  }
  backend "http" {
  }
}

provider "cloudflare" {
}

provider "routeros" {
  alias    = "dns"
  hosturl  = var.dns_routeros_hosturl
  username = var.dns_routeros_username
  password = var.dns_routeros_password
}

variable "dns_enabled" {
  description = "Manage the download aliases and delegation in both views."
  type        = bool
  default     = true
}

variable "dns_cloudflare_zone_id" {
  type    = string
  default = ""
}

variable "dns_routeros_hosturl" {
  type    = string
  default = ""
}

variable "dns_routeros_username" {
  type      = string
  default   = ""
  sensitive = true
}

variable "dns_routeros_password" {
  type      = string
  default   = ""
  sensitive = true
}

module "dns" {
  source             = "../../../projects/tf_modules/dns_records"
  document           = jsondecode(file("${path.module}/../dnsconfig.json"))
  cloudflare_zone_id = var.dns_cloudflare_zone_id
  enabled            = var.dns_enabled
  providers = {
    cloudflare = cloudflare
    routeros   = routeros.dns
  }
}

output "dns_import_addresses" {
  value = module.dns.import_addresses
}
