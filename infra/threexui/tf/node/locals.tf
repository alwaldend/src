locals {
  mullvad_relays = jsondecode(file("${path.module}/../../../../data/mullvad/relays.json"))
  mullvad_wireguard = {
    for relay in local.mullvad_relays.wireguard.relays : relay.hostname => relay
  }
}
