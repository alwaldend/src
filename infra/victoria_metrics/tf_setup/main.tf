terraform {
  required_providers {
    xenorchestra = {
      source  = "vatesfr/xenorchestra"
      version = "0.41.0"
    }
  }
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
  name = "src_infra_victoria_metrics"
}

locals {
  guest_address = jsondecode(file("${path.module}/../dns/dnsconfig.json")).records.upstream.A.address
  hostname      = "victoria-metrics"
  mac = join(":", concat(["02"], [
    for octet in range(5) : substr(sha256(local.hostname), octet * 2, 2)
  ]))
}

resource "xenorchestra_vm" "victoria_metrics" {
  name_label       = local.hostname
  name_description = "VictoriaMetrics single-node guest managed by infra/victoria_metrics"
  template         = data.xenorchestra_template.guest.id
  resource_set     = data.xenorchestra_resource_set.guest.id
  cpus             = 2
  memory_max       = 4 * 1024 * 1024 * 1024
  memory_min       = 4 * 1024 * 1024 * 1024
  auto_poweron     = true
  power_state      = "Running"
  tags             = ["victoria_metrics", "src_infra_victoria_metrics"]

  cloud_config = "#cloud-config\n${yamlencode(merge(
    jsondecode(file("${path.module}/../../cloud_init/xen_linux.json")),
    {
      hostname = local.hostname
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
        address         = "${local.guest_address}/24"
        gateway         = "192.168.10.1"
        dns_nameservers = ["192.168.10.1"]
      }]
    }]
  })

  network {
    network_id       = data.xenorchestra_network.guest.id
    mac_address      = local.mac
    expected_ip_cidr = "${local.guest_address}/32"
  }

  disk {
    sr_id      = data.xenorchestra_sr.guest.id
    name_label = "${local.hostname} boot and metrics"
    size       = 40 * 1024 * 1024 * 1024
  }

  timeouts {
    create = "20m"
  }
}

output "guest_address" {
  description = "Management IPv4 address of the VictoriaMetrics guest."
  value       = local.guest_address
}
