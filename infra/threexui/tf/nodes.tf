locals {
  nodes = {
    njalla1 = {
      address      = "njalla1.nodes.threexui.alwaldend.com"
      base_path    = split("/", var.njalla_xui_url)[length(split("/", var.njalla_xui_url)) - 1]
      api_token    = var.njalla_xui_token
      xui_username = var.njalla_xui_username
      xui_password = var.njalla_xui_password
      xui_url      = var.njalla_xui_url
      mullvad_key  = var.mullvad_keys["njalla1"]
    }
  }
}

provider "threexui" {
  alias    = "njalla1"
  endpoint = local.nodes["njalla1"].xui_url
  username = local.nodes["njalla1"].xui_username
  password = local.nodes["njalla1"].xui_password
}

module "node_njalla1" {
  source = "./node"
  providers = {
    threexui = threexui.njalla1
  }
  xui_base_path = local.nodes["njalla1"].base_path
  xui_address   = local.nodes["njalla1"].address
  xui_url       = local.nodes["njalla1"].xui_url
  xui_username  = local.nodes["njalla1"].xui_username
  xui_password  = local.nodes["njalla1"].xui_password
  mullvad_key   = local.nodes["njalla1"].mullvad_key
}

resource "threexui_node" "nodes" {
  for_each        = local.nodes
  name            = each.key
  remark          = each.value.address
  scheme          = "https"
  address         = each.value.address
  port            = 443
  base_path       = "/${each.value.base_path}/"
  api_token       = each.value.api_token
  enable          = true
  outbound_tag    = ""
  tls_verify_mode = "mtls"
}
