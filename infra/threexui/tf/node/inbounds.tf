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
  listen              = "127.0.0.1:40000"
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
      path = "/${var.xui_base_path}-inbound-40000"
    }
  }
}
