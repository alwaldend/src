terraform {
  required_providers {
    routeros = {
      source  = "terraform-routeros/routeros"
      version = "1.99.1"
    }
  }
}

variable "routeros_hosturl" {
  type        = string
  description = "RouterOS management endpoint supplied by the firewall credential injector."
}

variable "routeros_username" {
  type        = string
  description = "Existing non-DNS RouterOS username supplied by the scoped Vault read grant."
  sensitive   = true
}

variable "routeros_password" {
  type        = string
  description = "Existing non-DNS RouterOS password supplied by the scoped Vault read grant."
  sensitive   = true
}

provider "routeros" {
  hosturl  = var.routeros_hosturl
  username = var.routeros_username
  password = var.routeros_password
}

locals {
  printer_address = jsondecode(file("${path.module}/../dnsconfig.json")).records.printer1_a.A.address
}

data "routeros_ip_firewall" "forward_drop" {
  rules {
    filter = {
      chain    = "forward"
      action   = "drop"
      comment  = "drop forward"
      disabled = "false"
    }
  }
}

resource "routeros_ip_firewall_filter" "wifi_printer" {
  action        = "accept"
  chain         = "forward"
  comment       = "tf[infra/mikrotik/firewall]: Wi-Fi access to printer1"
  disabled      = false
  in_interface  = "bridge2"
  out_interface = "bridge1"
  src_address   = "192.168.2.0/24"
  dst_address   = local.printer_address
  place_before  = try(one(data.routeros_ip_firewall.forward_drop.rules[*].id), null)

  lifecycle {
    precondition {
      condition     = length(data.routeros_ip_firewall.forward_drop.rules) == 1
      error_message = "Expected exactly one enabled forward drop with comment 'drop forward'; inspect the live firewall before deploying."
    }
  }
}
