# resource "threexui_inbound_client" "freedom" {
#   for_each = local.admins_by_name
#   inbound_id  = threexui_inbound.freedom.id
#   email       = each.value.metadata.email
#   enable      = true
#   flow = "" # "xtls-rprx-vision"
# }


output "subs_freedom" {
  value = {
    #  for key, value in threexui_inbound_client.freedom : key => "${threexui_panel_subscription.settings.sub_uri}${value.sub_id}"
  }
}

resource "threexui_inbound" "freedom" {
  port                = 40000
  listen              = "127.0.0.1"
  protocol            = "vless"
  enable              = true
  remark              = "Freedom"
  share_addr_strategy = "node"
  vless_settings {
    decryption = "none"
  }
  stream_settings {
    network  = "xhttp"
    security = "none"
    xhttp_settings {
      path = "/inbound_${local.base_path_trimmed}_40000"
    }
  }
  sniffing {
    enabled       = true
    dest_override = ["http", "tls", "quic", "fakedns"]
  }
}

locals {
  mullvad_min = {
    for idx, name in var.mullvad_min : name => {
      wg   = local.mullvad_wireguard[name]
      port = 40010 + idx
      name = name
    }
  }
  http_proxy_accounts = keys(var.http_proxies)
  http_proxy = {
    for idx, name in var.mullvad_min : name => {
      port = 40080 + idx
      name = name
      user = local.http_proxy_accounts[idx % length(local.http_proxy_accounts)]
      pass = var.http_proxies[local.http_proxy_accounts[idx % length(local.http_proxy_accounts)]]
    }
  }
  mikrotik_parent_proxy = var.mikrotik_parent_proxy == null ? {} : {
    (var.mikrotik_parent_proxy.outbound) = var.mikrotik_parent_proxy
  }
}

# resource "threexui_inbound_client" "mullvad_min" {
#   for_each = merge([
#     for entity_name, entity in local.admins_by_name : {
#       for inbound_name, inbound in threexui_inbound.mullvad_min : "${entity_name}_${inbound_name}" => {
#         inbound = inbound
#         entity  = entity
#       }
#     }
#   ]...)
#   inbound_id = each.value.inbound.id
#   email      = each.value.entity.metadata.email
#   enable     = true
#   flow       = "xtls-rprx-vision"
# }

resource "threexui_inbound" "mullvad_min" {
  for_each            = local.mullvad_min
  port                = each.value.port
  listen              = "127.0.0.1"
  protocol            = "vless"
  enable              = true
  remark              = "mullvad_min | wg | ${each.key}"
  share_addr_strategy = "node"
  vless_settings {
    decryption = "none"
  }
  stream_settings {
    network  = "xhttp"
    security = "none"
    xhttp_settings {
      path = "/inbound_${local.base_path_trimmed}_${each.value.port}"
    }
  }
  sniffing {
    enabled       = true
    dest_override = ["http", "tls", "quic", "fakedns"]
  }
}

resource "threexui_inbound" "http_proxy" {
  for_each            = local.http_proxy
  port                = each.value.port
  listen              = "127.0.0.1"
  protocol            = "http"
  enable              = true
  remark              = "http_proxy | ${each.key}"
  share_addr_strategy = "node"

  # Xray enables authentication when accounts are present. The HTTP inbound
  # schema has no separate auth field.
  http_settings {
    allow_transparent = false
    account {
      user = each.value.user
      pass = each.value.pass
    }
  }
}

resource "threexui_inbound" "mikrotik_parent_proxy" {
  for_each            = local.mikrotik_parent_proxy
  port                = each.value.port
  listen              = "0.0.0.0"
  protocol            = "http"
  enable              = true
  remark              = "http_proxy | mikrotik parent | ${each.key}"
  share_addr_strategy = "node"

  # No accounts means that the proxy is intentionally unauthenticated.
  http_settings {
    allow_transparent = false
  }

  lifecycle {
    precondition {
      condition     = contains(var.mullvad_min, each.key)
      error_message = "The MikroTik parent proxy outbound must be present in mullvad_min."
    }
  }
}

resource "threexui_inbound" "mullvad_min_lb" {
  port                = 40050
  listen              = "127.0.0.1"
  protocol            = "vless"
  enable              = true
  remark              = "mullvad_min | lb"
  share_addr_strategy = "node"
  vless_settings {
    decryption = "none"
  }
  stream_settings {
    network  = "xhttp"
    security = "none"
    xhttp_settings {
      path = "/inbound_${local.base_path_trimmed}_40050"
    }
  }
  sniffing {
    enabled       = true
    dest_override = ["http", "tls", "quic", "fakedns"]
  }
}
