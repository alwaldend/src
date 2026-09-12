terraform {
  required_version = ">= 1.9"
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "= 5.22.0"
    }
    routeros = {
      source  = "terraform-routeros/routeros"
      version = "= 1.99.1"
    }
  }
}

variable "document" {
  type        = any
  description = "Decoded canonical dnsconfig.json consumed by the provider-free normalizer."
}

variable "zone" {
  type        = string
  default     = "alwaldend.com"
  description = "DNS zone containing the declared names."
}

variable "cloudflare_zone_id" {
  type        = string
  default     = null
  description = "Optional Cloudflare zone identifier; enabled global records otherwise resolve the zone by name."
}

variable "enabled" {
  type        = bool
  default     = false
  description = "Enable provider resource ownership only after the coordinated DNSControl exclusion and import preparation. Keep true after adoption."
}

module "global" {
  source             = "./global"
  document           = var.document
  zone               = var.zone
  cloudflare_zone_id = var.cloudflare_zone_id
  enabled            = var.enabled
}

locals {
  dc1_records = { for key, record in module.global.normalized_records : key => record if record.view == "dc1" }
}

module "dc1" {
  source             = "./dc1"
  count              = var.enabled && length(local.dc1_records) > 0 ? 1 : 0
  normalized_records = local.dc1_records
}

output "normalized_records" {
  description = "Complete declared records, including while provider ownership is disabled."
  value       = module.global.normalized_records
}

output "import_addresses" {
  description = "Resource addresses relative to this module, including records staged while disabled; provider IDs require a live inventory."
  value = merge(
    { for key, address in module.global.import_addresses : key => "module.global.${address}" },
    { for key, record in module.global.normalized_records : key => format(
      "module.dc1[0].routeros_ip_dns_record.records[%s]", jsonencode(key),
    ) if record.view == "dc1" },
  )
}
