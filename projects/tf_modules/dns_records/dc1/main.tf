terraform {
  required_providers {
    routeros = {
      source  = "terraform-routeros/routeros"
      version = "= 1.99.1"
    }
  }
}

variable "normalized_records" {
  description = "Validated canonical dc1 records selected by the owning parent module."
  type = map(object({
    name     = string
    type     = string
    value    = string
    priority = number
    ttl      = number
  }))
}

resource "routeros_ip_dns_record" "records" {
  for_each = var.normalized_records

  name            = each.value.name
  type            = each.value.type
  address         = contains(["A", "AAAA"], each.value.type) ? each.value.value : null
  cname           = each.value.type == "CNAME" ? each.value.value : null
  mx_exchange     = each.value.type == "MX" ? each.value.value : null
  mx_preference   = each.value.priority
  ns              = each.value.type == "NS" ? each.value.value : null
  text            = each.value.type == "TXT" ? each.value.value : null
  ttl             = "${each.value.ttl}s"
  disabled        = false
  match_subdomain = false
  comment         = ""
}

output "records" {
  description = "Bound provider records for the combined module's ownership checks."
  value       = routeros_ip_dns_record.records
}
