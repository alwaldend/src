resource "threexui_xray_balancers" "config" {
  balancer {
    tag = "out-mullvad-min-lb"
    selector = [
      for key, _ in threexui_inbound.mullvad_min : "out-mullvad-min-${key}"
    ]
    fallback_tag = "blocked"
    strategy {
      type = "leastLoad"
    }
  }
}
