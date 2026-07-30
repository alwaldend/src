resource "routeros_interface_wireguard" "vpc" {
  name        = "ingress-vpc"
  comment     = "ingress-vpc (infra/ingress/tf)"
  listen_port = "13231"
  mtu         = "1420"
  private_key = var.wg_private_keys.router
}

resource "routeros_interface_wireguard_peer" "vpc" {
  for_each             = local.vpc
  name                 = "${routeros_interface_wireguard.vpc.name}-${each.key}"
  comment              = each.key
  interface            = routeros_interface_wireguard.vpc.name
  public_key           = var.wg_public_keys[each.key]
  preshared_key        = var.wg_preshared_keys[each.key]
  endpoint_address     = yandex_vpc_address.vpc[each.key].external_ipv4_address[0].address
  endpoint_port        = "51820"
  persistent_keepalive = "5s"
  allowed_address = [
    "0.0.0.0/0",
    "::/0"
  ]
}
