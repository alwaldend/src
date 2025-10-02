terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5"
    }
  }
}

provider "cloudflare" {
}

locals {
    data = yamldecode(file("${path.module}/../../data/dns.yaml"))
    records_alwaldend_com = { for record in local.data.records["alwaldend.com"]: "${record.name}-${record.type}-${record.content}" => merge(record, { zone_id = local.data.cloudflare_zones["alwaldend.com"].id }) }
}

resource "cloudflare_dns_record" "alwaldend_com" {
    for_each = local.records_alwaldend_com
    zone_id = each.value.zone_id
    comment = try(each.value.comment, "")
    content = each.value.content
    name = each.value.name
    tags = try(each.value.tags, [])
    priority = try(each.value.priority, 0)
    ttl = each.value.ttl
    type = each.value.type
}
