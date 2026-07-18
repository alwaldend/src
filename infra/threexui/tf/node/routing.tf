resource "threexui_xray_routing" "config" {
  domain_strategy = "AsIs"
  domain_matcher  = "hybrid"

  rule {
    domain       = ["geosite:category-ads"]
    outbound_tag = "blocked"
  }

  rule {
    ip           = ["geoip:private"]
    outbound_tag = "blocked"
  }

  rule {
    inbound_tag  = [threexui_inbound.freedom.tag]
    outbound_tag = "direct"
  }
}
