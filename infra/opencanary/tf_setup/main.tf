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

data "xenorchestra_template" "guest" {
  name_label = "fedora44-cloud.templates.xcp-ng.alwaldend.com"
}

data "xenorchestra_sr" "guest" {
  name_label = "Local storage"
}

data "xenorchestra_network" "guest" {
  name_label = "Pool-wide network 1"
}

data "xenorchestra_resource_set" "guest" {
  name = "src_infra_opencanary"
}

locals {
  record   = jsondecode(file("${path.module}/../dns/dnsconfig.json")).records.opencanary.A
  hostname = "${local.record.name}.alwaldend.com"
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
    {
      hostname = local.record.name
      fqdn     = local.hostname
      runcmd   = [["systemctl", "enable", "--now", "xe-linux-distribution.service"]]
    },
  ))}"

  cloud_network_config = yamlencode({
    version = 1
    config = [{
      type        = "physical"
      name        = "enX0"
      mac_address = local.mac
      subnets = [{
        type            = "static"
        address         = "${local.record.address}/24"
        gateway         = "192.168.10.1"
        dns_nameservers = ["192.168.10.1"]
      }]
    }]
  })

  network {
    network_id       = data.xenorchestra_network.guest.id
    mac_address      = local.mac
    expected_ip_cidr = "${local.record.address}/32"
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
  description = "Static IPv4 address of the guest."
  value       = local.record.address
}
