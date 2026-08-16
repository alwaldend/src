locals {
  nodes = {
    njalla1 = {
      address = "njalla1.nodes.threexui.alwaldend.com"
      mullvad_min = [
        "nl-ams-wg-008",
        "us-sjc-wg-003",
        "sg-sin-wg-003",
        "de-dus-wg-003",
        "fr-mrs-wg-002",
      ]
    }
    yc1 = {
      address = "yc1.nodes.threexui.alwaldend.com"
      mullvad_min = [
        "us-sjc-wg-001",
        "dk-cph-wg-102",
        "ca-mtr-wg-301",
        "jp-osa-wg-102",
        "nz-akl-wg-302",
      ]
    }
  }
}

provider "threexui" {
  alias    = "njalla1"
  endpoint = "https://${local.nodes.njalla1.address}${var.node_xui_base_path}"
  username = var.node_xui_username
  password = var.node_xui_password
}

provider "threexui" {
  alias    = "yc1"
  endpoint = "https://${local.nodes.yc1.address}${var.node_xui_base_path}"
  username = var.node_xui_username
  password = var.node_xui_password
}

module "node_njalla1" {
  source = "./node"
  providers = {
    threexui = threexui.njalla1
  }
  xui_base_path = var.node_xui_base_path
  xui_address   = local.nodes.njalla1.address
  xui_url       = "https://${local.nodes.njalla1.address}${var.node_xui_base_path}"
  xui_username  = var.node_xui_username
  xui_password  = var.node_xui_password
  mullvad_key   = var.mullvad_keys["njalla1"]
  mullvad_min   = local.nodes.njalla1.mullvad_min
  http_proxies  = var.http_proxies
}

module "node_yc1" {
  source = "./node"
  providers = {
    threexui = threexui.yc1
  }
  xui_base_path = var.node_xui_base_path
  xui_address   = local.nodes.yc1.address
  xui_url       = "https://${local.nodes.yc1.address}${var.node_xui_base_path}"
  xui_username  = var.node_xui_username
  xui_password  = var.node_xui_password
  mullvad_key   = var.mullvad_keys["yc1"]
  mullvad_min   = local.nodes.yc1.mullvad_min
  http_proxies  = var.http_proxies
}

resource "threexui_node" "nodes" {
  for_each        = local.nodes
  name            = each.key
  remark          = each.value.address
  scheme          = "https"
  address         = each.value.address
  port            = 443
  base_path       = "${var.node_xui_base_path}/"
  api_token       = var.node_xui_tokens[each.key]
  enable          = true
  outbound_tag    = ""
  tls_verify_mode = "mtls"
}
