mock_provider "cloudflare" {
  mock_data "cloudflare_zone" {
    defaults = {
      id   = "fedcba9876543210fedcba9876543210"
      name = "alwaldend.com"
    }
  }
}

variables {
  cloudflare_zone_id = "0123456789abcdef0123456789abcdef"
  document = {
    records = {
      host   = { A = { name = "host", address = "192.0.2.1" }, AAAA = { name = "host", address = "2001:db8::1" }, dsp = ["global"] }
      alias  = { CNAME = { name = "alias", target = "host" }, dsp = ["global"] }
      ns     = { NS = { name = "child", address = "ns.example.net." }, dsp = ["global"] }
      mx     = { MX = { name = "@", target = "mail.example.net.", priority = 20, ttl = 10800 }, dsp = ["global"] }
      txt    = { TXT = { name = "@", content = "v=spf1 -all", ttl = 300 }, dsp = ["global"] }
      second = { A = { name = "host", address = "192.0.2.2" }, dsp = ["global"] }
    }
  }
}

run "staged_without_provider_ownership" {
  command = plan
  variables { cloudflare_zone_id = null }
  assert {
    condition     = length(cloudflare_dns_record.records) == 0 && length(data.cloudflare_zone.selected) == 0 && length(output.normalized_records) == 7
    error_message = "Global staging without a zone ID must preserve declarations without looking up the zone or creating provider records."
  }
  assert {
    condition     = output.import_addresses["host/A/global"] == "cloudflare_dns_record.records[\"host/A/global\"]"
    error_message = "Global import addresses must be relative to the selected entrypoint."
  }
}

run "global_only_provider_mapping" {
  command = apply
  variables { enabled = true }
  assert {
    condition     = length(cloudflare_dns_record.records) == 7 && alltrue([for record in cloudflare_dns_record.records : record.proxied == false])
    error_message = "Global owners must create only their declared unproxied records without a RouterOS provider."
  }
  assert {
    condition     = length(data.cloudflare_zone.selected) == 0 && alltrue([for record in cloudflare_dns_record.records : record.zone_id == var.cloudflare_zone_id])
    error_message = "An explicit zone ID must be used unchanged without a zone lookup."
  }
  assert {
    condition     = cloudflare_dns_record.records["host/A/global"].content == "192.0.2.1" && cloudflare_dns_record.records["host/AAAA/global"].content == "2001:db8::1" && cloudflare_dns_record.records["alias/CNAME/global"].content == "host.alwaldend.com" && cloudflare_dns_record.records["ns/NS/global"].content == "ns.example.net" && cloudflare_dns_record.records["txt/TXT/global"].content == "v=spf1 -all"
    error_message = "Every Cloudflare scalar type must use its normalized content."
  }
  assert {
    condition     = cloudflare_dns_record.records["mx/MX/global"].content == "mail.example.net" && cloudflare_dns_record.records["mx/MX/global"].priority == 20 && cloudflare_dns_record.records["mx/MX/global"].ttl == 10800
    error_message = "Cloudflare must preserve MX priority and explicit TTL."
  }
}

run "resolve_missing_zone_id" {
  command = plan
  variables {
    enabled            = true
    cloudflare_zone_id = null
    zone               = "AlWaLdEnD.CoM."
  }
  assert {
    condition     = length(data.cloudflare_zone.selected) == 1 && data.cloudflare_zone.selected[0].filter.name == "alwaldend.com" && alltrue([for record in cloudflare_dns_record.records : record.zone_id == "fedcba9876543210fedcba9876543210"])
    error_message = "Enabled global records without a zone ID must use the matching zone's discovered ID."
  }
}

run "resolve_empty_zone_id" {
  command = plan
  variables {
    enabled            = true
    cloudflare_zone_id = ""
  }
  assert {
    condition     = length(data.cloudflare_zone.selected) == 1 && alltrue([for record in cloudflare_dns_record.records : record.zone_id == "fedcba9876543210fedcba9876543210"])
    error_message = "An empty optional zone ID from the injector must use zone-name discovery."
  }
}

run "dc1_only_does_not_lookup_zone" {
  command = plan
  variables {
    enabled            = true
    cloudflare_zone_id = null
    document = {
      records = {
        host = { A = { name = "host", address = "192.0.2.1" }, dsp = ["dc1"] }
      }
    }
  }
  assert {
    condition     = length(data.cloudflare_zone.selected) == 0 && length(cloudflare_dns_record.records) == 0 && length(output.normalized_records) == 1
    error_message = "An owner without global records must preserve normalization without looking up a Cloudflare zone."
  }
}

run "updates_preserve_ids_and_other_values" {
  command = apply
  variables {
    enabled = true
    document = {
      records = {
        host   = { A = { name = "host", address = "192.0.2.99" }, dsp = ["global"] }
        second = { A = { name = "host", address = "192.0.2.2" }, dsp = ["global"] }
      }
    }
  }
  assert {
    condition     = length(cloudflare_dns_record.records) == 2 && cloudflare_dns_record.records["host/A/global"].content == "192.0.2.99" && cloudflare_dns_record.records["second/A/global"].content == "192.0.2.2" && output.record_ids["host/A/global"] == run.global_only_provider_mapping.record_ids["host/A/global"] && output.record_ids["second/A/global"] == run.global_only_provider_mapping.record_ids["second/A/global"]
    error_message = "A value change and removal of other owned records must preserve both retained provider IDs."
  }
}
