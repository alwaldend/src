resource "routeros_interface_wireguard" "vpc" {
  name        = "ingress-vpc"
  comment     = "tf[infra/ingress/tf]"
  listen_port = "13231"
  mtu         = "1420"
  private_key = var.wg_private_keys.router
}

locals {
  inventory = yamldecode(file("${path.module}/../ansible/inventory.yaml"))
}

resource "routeros_ip_address" "vpc" {
  address   = "10.10.0.0/24"
  interface = routeros_interface_wireguard.vpc.name
  comment   = "tf[infra/ingress/tf]"
}

resource "routeros_interface_list_member" "vpc" {
  for_each  = { LAN = {}, accept-forward-LAN = {}, accept-input-ICMP = {} }
  interface = routeros_interface_wireguard.vpc.name
  list      = each.key
  comment   = "tf[infra/ingress/tf]"
}

resource "routeros_interface_wireguard_peer" "vpc" {
  for_each             = local.vpc
  name                 = "${routeros_interface_wireguard.vpc.name}-${each.key}"
  comment              = each.key
  interface            = routeros_interface_wireguard.vpc.name
  private_key          = ""
  public_key           = var.wg_public_keys[each.key]
  preshared_key        = var.wg_preshared_keys[each.key]
  endpoint_address     = yandex_vpc_address.vpc[each.key].external_ipv4_address[0].address
  endpoint_port        = "51820"
  persistent_keepalive = "5s"
  allowed_address = [
    local.inventory.ingress.hosts["${each.key}.ingress.alwaldend.com"].wg_ip,
  ]
}
