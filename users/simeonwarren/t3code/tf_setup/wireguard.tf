resource "routeros_interface_wireguard" "t3code" {
  name        = "t3code-vpc"
  comment     = "tf[users/simeonwarren/t3code/tf_setup]"
  listen_port = "13233"
  mtu         = "1420"
  private_key = var.wg_private_keys.router
}

locals {
  inventory = yamldecode(file("${path.module}/../ansible/inventory.yaml"))
}

resource "routeros_ip_address" "t3code" {
  address   = "10.30.1.0/24"
  interface = routeros_interface_wireguard.t3code.name
  comment   = "tf[users/simeonwarren/t3code/tf_setup]"
}

resource "routeros_interface_list_member" "t3code" {
  for_each  = { LAN = {}, accept-forward-LAN = {}, accept-input-ICMP = {} }
  interface = routeros_interface_wireguard.t3code.name
  list      = each.key
  comment   = "tf[users/simeonwarren/t3code/tf_setup]"
}

resource "routeros_interface_wireguard_peer" "t3code" {
  for_each             = local.yc_vms
  name                 = "${routeros_interface_wireguard.t3code.name}-${each.key}"
  comment              = each.key
  interface            = routeros_interface_wireguard.t3code.name
  private_key          = ""
  public_key           = var.wg_public_keys[each.key]
  preshared_key        = var.wg_preshared_keys[each.key]
  endpoint_address     = yandex_compute_instance.t3code[each.key].network_interface[0].nat_ip_address
  endpoint_port        = "51820"
  persistent_keepalive = "5s"
  allowed_address = [
    local.inventory.t3code.hosts["${each.key}.t3code.simeonwarren.users.alwaldend.com"].wg_ip,
    "192.168.1.0/24",
  ]
}
