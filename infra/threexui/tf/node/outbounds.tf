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

  dynamic "outbound" {
    for_each = local.mullvad_min
    content {
      tag      = "out-mullvad-min-${outbound.key}"
      protocol = "wireguard"
      wireguard_settings {
        secret_key      = var.mullvad_key.private_key
        address         = split(",", var.mullvad_key.address)
        mtu             = 1420
        domain_strategy = "ForceIPv6v4"
        peer {
          public_key  = outbound.value.wg.public_key
          endpoint    = "${outbound.value.wg.ipv4_addr_in}:51820"
          allowed_ips = ["0.0.0.0/0", "::/0"]
        }
      }
    }
  }
}
