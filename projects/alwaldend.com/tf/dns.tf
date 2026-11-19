module "dns" {
  source = "../../tf_modules/dns_records/global"

  document           = jsondecode(file("${path.module}/../dnsconfig.json"))
  cloudflare_zone_id = var.dns_cloudflare_zone_id
  enabled            = var.dns_enabled
}

variable "dns_cloudflare_import_ids" {
  type        = map(string)
  description = "Existing Cloudflare zone/record IDs keyed by normalized DNS record key; supplied only during reviewed adoption."
  default     = {}
}

import {
  for_each = var.dns_cloudflare_import_ids
  to       = module.dns.cloudflare_dns_record.records[each.key]
  id       = each.value
}
