locals {
  # YC1 assigns HTTP ports from 40080 in mullvad_min order. Port 40080 is
  # the existing us-sjc-wg-001 HTTP inbound.
  mikrotik_http_proxy = {
    listen_port = 8080
    parent_host = "yc1.nodes.threexui.alwaldend.com"
    parent_port = 40080
  }
  mikrotik_http_proxy_script = <<-ROS
    :local parentProxy [:resolve "${local.mikrotik_http_proxy.parent_host}"]
    /ip/proxy/set enabled=yes port=${local.mikrotik_http_proxy.listen_port} parent-proxy=$parentProxy parent-proxy-port=${local.mikrotik_http_proxy.parent_port}
  ROS
}

data "routeros_ip_firewall" "drop_input" {
  rules {
    filter = {
      action  = "drop"
      chain   = "input"
      comment = "drop input"
    }
  }
}

resource "routeros_system_script" "http_proxy_parent" {
  name           = "t3code-http-proxy-parent"
  comment        = "tf[users/simeonwarren/t3code/tf_setup]"
  source         = local.mikrotik_http_proxy_script
  launch_trigger = sha256(local.mikrotik_http_proxy_script)
  policy         = ["read", "write"]
}

resource "routeros_system_scheduler" "http_proxy_parent" {
  name     = "t3code-http-proxy-parent"
  comment  = "tf[users/simeonwarren/t3code/tf_setup]"
  on_event = routeros_system_script.http_proxy_parent.name
  interval = "5m"
  policy   = ["read", "write"]
}

resource "routeros_ip_firewall_filter" "http_proxy" {
  action       = "accept"
  chain        = "input"
  comment      = "tf[users/simeonwarren/t3code/tf_setup] HTTP proxy"
  dst_port     = tostring(local.mikrotik_http_proxy.listen_port)
  in_interface = routeros_interface_wireguard.t3code.name
  place_before = data.routeros_ip_firewall.drop_input.rules[0].id
  protocol     = "tcp"
}
