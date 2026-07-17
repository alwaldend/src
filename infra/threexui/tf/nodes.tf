locals {
  nodes = {
    njalla1 = {
      address   = "njalla1.nodes.threexui.alwaldend.com"
      url_split = split("/", var.njalla_xui_url)
      api_token = var.njalla_xui_token
    }
  }
}

resource "threexui_node" "nodes" {
  for_each        = local.nodes
  name            = each.key
  remark          = each.value.address
  scheme          = "https"
  address         = each.value.address
  port            = 443
  base_path       = join("/", ["", each.value.url_split[length(each.value.url_split) - 1], ""])
  api_token       = each.value.api_token
  enable          = true
  outbound_tag    = ""
  tls_verify_mode = "mtls"
}
