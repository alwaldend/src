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
}

variable "dns_routeros_hosturl" {
  type        = string
  description = "RouterOS management endpoint supplied by the DNS credential injector."
  default     = ""
}

variable "dns_routeros_username" {
  type        = string
  description = "RouterOS DNS username supplied by the DNS credential injector."
  sensitive   = true
  default     = ""
}

variable "dns_routeros_password" {
  type        = string
  description = "RouterOS DNS password supplied by the DNS credential injector."
  sensitive   = true
  default     = ""
}

provider "routeros" {
  alias    = "dns"
  hosturl  = var.dns_routeros_hosturl
  username = var.dns_routeros_username
  password = var.dns_routeros_password
}

module "dns" {
  source = "../../../projects/tf_modules/dns_records"

  providers = {
    cloudflare = cloudflare
    routeros   = routeros.dns
  }

  document           = jsondecode(file("${path.module}/dnsconfig.json"))
  cloudflare_zone_id = var.dns_cloudflare_zone_id
  enabled            = var.dns_enabled
}
