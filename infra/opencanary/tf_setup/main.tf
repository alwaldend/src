terraform {
  required_providers {
    xenorchestra = {
      source  = "vatesfr/xenorchestra"
      version = "0.41.0"
    }
  }

  backend "http" {}
}

provider "xenorchestra" {
  url = "wss://xoa.xcp-ng.alwaldend.com"
}

variable "network_name" {
  type        = string
  description = "Name of the isolated honeypot network assigned to the OpenCanary resource set."
}

data "xenorchestra_template" "guest" {
  name_label = "fedora44-cloud.templates.xcp-ng.alwaldend.com"
}

data "xenorchestra_sr" "guest" {
  name_label = "Local storage"
}

data "xenorchestra_network" "guest" {
  name_label = var.network_name
}

data "xenorchestra_resource_set" "guest" {
  name = "src_infra_opencanary"
}

locals {
  hostname = "opencanary.alwaldend.com"
  mac = join(":", concat(["02"], [
    for octet in range(5) : substr(sha256(local.hostname), octet * 2, 2)
  ]))
}

resource "xenorchestra_vm" "opencanary" {
  name_label       = local.hostname
  name_description = "OpenCanary honeypot managed by infra/opencanary"
  template         = data.xenorchestra_template.guest.id
  resource_set     = data.xenorchestra_resource_set.guest.id
  cpus             = 2
  memory_max       = 2 * 1024 * 1024 * 1024
  memory_min       = 2 * 1024 * 1024 * 1024
  auto_poweron     = true
  power_state      = "Running"
  tags             = ["opencanary", "src_infra_opencanary"]

  cloud_config = "#cloud-config\n${yamlencode(merge(
    jsondecode(file("${path.module}/../../cloud_init/xen_linux.json")),
    { hostname = "opencanary", fqdn = local.hostname },
  ))}"

  network {
    network_id       = data.xenorchestra_network.guest.id
    mac_address      = local.mac
    expected_ip_cidr = "0.0.0.0/0"
  }

  disk {
    sr_id      = data.xenorchestra_sr.guest.id
    name_label = "${local.hostname} boot"
    size       = 20 * 1024 * 1024 * 1024
  }

  timeouts {
    create = "20m"
  }
}

output "address" {
  description = "Guest IPv4 address reported by Xen Orchestra, for OPENCANARY_HOST."
  value       = try(xenorchestra_vm.opencanary.network[0].ipv4_addresses[0], null)
}

output "mac_address" {
  description = "Stable guest MAC address for an optional DHCP reservation."
  value       = local.mac
}
