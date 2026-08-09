resource "routeros_interface_wireguard" "hermes" {
  name        = "hermes-vpc"
  comment     = "tf[users/simeonwarren/hermes/tf_setup]"
  listen_port = "13232"
  mtu         = "1420"
  private_key = var.wg_private_keys.router
}

locals {
  inventory = yamldecode(file("${path.module}/../ansible/inventory.yaml"))
}

resource "routeros_ip_address" "hermes" {
  address   = "10.20.0.0/24"
  interface = routeros_interface_wireguard.hermes.name
  comment   = "tf[users/simeonwarren/hermes/tf_setup]"
}

resource "routeros_interface_list_member" "hermes" {
  for_each  = { LAN = {}, accept-forward-LAN = {}, accept-input-ICMP = {} }
  interface = routeros_interface_wireguard.hermes.name
  list      = each.key
  comment   = "tf[users/simeonwarren/hermes/tf_setup]"
}

resource "routeros_interface_wireguard_peer" "hermes" {
  for_each             = local.yc_vms
  name                 = "${routeros_interface_wireguard.hermes.name}-${each.key}"
  comment              = each.key
  interface            = routeros_interface_wireguard.hermes.name
  private_key          = ""
  public_key           = var.wg_public_keys[each.key]
  preshared_key        = var.wg_preshared_keys[each.key]
  endpoint_address     = yandex_compute_instance.hermes[each.key].network_interface[0].nat_ip_address
  endpoint_port        = "51820"
  persistent_keepalive = "5s"
  allowed_address = [
    local.inventory.hermes.hosts["${each.key}.hermes.simeonwarren.users.alwaldend.com"].wg_ip,
  ]
}
