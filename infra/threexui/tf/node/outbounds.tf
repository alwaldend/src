resource "threexui_xray_outbounds" "config" {
  outbound {
    tag      = "direct"
    protocol = "freedom"

    freedom_settings {
      domain_strategy = "AsIs"
    }
  }

  outbound {
    tag      = "blocked"
    protocol = "blackhole"

    blackhole_settings {
      response_type = "none"
    }
  }
}
