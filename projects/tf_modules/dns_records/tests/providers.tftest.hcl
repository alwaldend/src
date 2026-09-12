mock_provider "cloudflare" {}
mock_provider "routeros" {}

variables {
  cloudflare_zone_id = "0123456789abcdef0123456789abcdef"
  document = {
    records = {
      host = {
        A    = { name = "host", address = "192.0.2.1" }
        AAAA = { name = "host", address = "2001:db8::1" }
        dsp  = ["all", "dc1"]
      }
      alias  = { CNAME = { name = "alias", target = "host" }, dsp = ["all"] }
      ns     = { NS = { name = "child", address = "ns.example.net." }, dsp = ["all"] }
      mx     = { MX = { name = "@", target = "mail.example.net.", priority = 20, ttl = 10800 }, dsp = ["all"] }
      txt    = { TXT = { name = "@", content = "v=spf1 -all", ttl = 300 }, dsp = ["all"] }
      second = { A = { name = "host", address = "192.0.2.2" }, dsp = ["global"] }
    }
  }
}

run "staged_without_provider_ownership" {
  command = plan
  assert {
    condition     = length(module.global.record_ids) == 0 && length(module.dc1) == 0 && length(output.normalized_records) == 13
    error_message = "The default-disabled module must preserve declarations without creating provider records."
  }
  assert {
    condition     = output.import_addresses["host/A/global"] == "module.global.cloudflare_dns_record.records[\"host/A/global\"]" && output.import_addresses["host/A/dc1"] == "module.dc1[0].routeros_ip_dns_record.records[\"host/A/dc1\"]"
    error_message = "Import guidance must use the actual stable resource addresses even while staged."
  }
}

run "enabled_provider_mapping" {
  command = apply
  variables { enabled = true }
  assert {
    condition     = length(module.global.record_ids) == 7 && length(module.dc1[0].records) == 6
    error_message = "Only the selected view's distinct scalar records may be created."
  }
  assert {
    condition     = module.dc1[0].records["host/A/dc1"].address == "192.0.2.1" && module.dc1[0].records["host/AAAA/dc1"].address == "2001:db8::1" && module.dc1[0].records["alias/CNAME/dc1"].cname == "host.alwaldend.com" && module.dc1[0].records["ns/NS/dc1"].ns == "ns.example.net" && module.dc1[0].records["mx/MX/dc1"].mx_exchange == "mail.example.net" && module.dc1[0].records["mx/MX/dc1"].mx_preference == 20 && module.dc1[0].records["txt/TXT/dc1"].text == "v=spf1 -all"
    error_message = "RouterOS must map each type to its scalar provider field."
  }
  assert {
    condition     = module.dc1[0].records["host/A/dc1"].ttl == "300s" && alltrue([for record in module.dc1[0].records : !record.disabled && !record.match_subdomain && record.comment == ""])
    error_message = "RouterOS must preserve exact-name enabled records and the effective TTL."
  }
}

run "update_and_delete_only_owned_records" {
  command = apply
  variables {
    enabled = true
    document = {
      records = {
        host = {
          A   = { name = "host", address = "192.0.2.99" }
          dsp = ["global"]
        }
        second = { A = { name = "host", address = "192.0.2.2" }, dsp = ["global"] }
      }
    }
  }
  assert {
    condition     = length(module.global.record_ids) == 2 && length(module.dc1) == 0 && output.normalized_records["host/A/global"].value == "192.0.2.99" && output.normalized_records["second/A/global"].value == "192.0.2.2"
    error_message = "Changing one value and removing owned declarations must retain the other scalar record."
  }
}
