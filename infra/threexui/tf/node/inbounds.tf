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
