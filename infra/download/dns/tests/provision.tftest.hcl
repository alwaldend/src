mock_provider "cloudflare" {}
mock_provider "routeros" { alias = "dns" }

variables {
  dns_cloudflare_zone_id = "0123456789abcdef0123456789abcdef"
}

# Exercise both provider mappings with the real canonical declaration.
run "split_horizon_records" {
  command = plan
}
