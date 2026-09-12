terraform {
  required_version = ">= 1.9"
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "= 5.22.0"
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
  description = "Enable provider resource ownership after coordinated DNSControl exclusion and import preparation. Keep true after adoption."
}

module "normalize" {
  source   = "../normalize"
  document = var.document
  zone     = var.zone
}

locals {
  global_records   = { for key, record in module.normalize.normalized_records : key => record if var.enabled && record.view == "global" }
  explicit_zone_id = var.cloudflare_zone_id != null && var.cloudflare_zone_id != ""
  zone_id          = local.explicit_zone_id ? var.cloudflare_zone_id : try(data.cloudflare_zone.selected[0].id, null)
}

data "cloudflare_zone" "selected" {
  count = length(local.global_records) > 0 && !local.explicit_zone_id ? 1 : 0
  filter = {
    name = module.normalize.zone
  }
}

resource "cloudflare_dns_record" "records" {
  for_each = local.global_records

  zone_id  = local.zone_id
  name     = each.value.name
  type     = each.value.type
  content  = each.value.value
  priority = each.value.priority
  ttl      = each.value.ttl
  proxied  = each.value.proxied

  lifecycle {
    precondition {
      condition     = local.zone_id != null && local.zone_id != ""
      error_message = "Enabled global DNS records require an explicit zone ID or a unique zone-name lookup result."
    }
  }
}

output "normalized_records" {
  description = "Complete canonical declarations for all views, including while provider ownership is disabled."
  value       = module.normalize.normalized_records
}

output "import_addresses" {
  description = "Global resource addresses relative to this module, including while staged; live record IDs are caller-owned."
  value = { for key, record in module.normalize.normalized_records : key => format(
    "cloudflare_dns_record.records[%s]", jsonencode(key),
  ) if record.view == "global" }
}

output "record_ids" {
  description = "Provider IDs for the global records currently bound to this module's state."
  value       = { for key, record in cloudflare_dns_record.records : key => record.id }
}
