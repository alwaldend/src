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
  name = "src_infra_mlflow"
}

locals {
  record        = jsondecode(file("${path.module}/../dns/dnsconfig.json")).records.local.A
  guest_address = local.record.address
  data_disks    = jsondecode(file("${path.module}/../storage.json"))
  hostname      = split(".", local.record.name)[0]
  mac = join(":", concat(["02"], [
    for octet in range(5) : substr(sha256(local.hostname), octet * 2, 2)
  ]))
}

resource "xenorchestra_vm" "mlflow" {
  name_label       = local.hostname
  name_description = "MLflow guest managed by infra/mlflow"
  template         = data.xenorchestra_template.guest.id
  resource_set     = data.xenorchestra_resource_set.guest.id
  cpus             = 2
  memory_max       = 4 * 1024 * 1024 * 1024
  memory_min       = 4 * 1024 * 1024 * 1024
  auto_poweron     = true
  power_state      = "Running"
  tags             = ["mlflow", "src_infra_mlflow"]

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
    name_label = "${local.hostname} boot"
    size       = 60 * 1024 * 1024 * 1024
  }

  # The template boot disk occupies xvda. Keep data disks in Xen device order.
  dynamic "disk" {
    for_each = local.data_disks
    content {
      sr_id      = data.xenorchestra_sr.guest.id
      name_label = "${local.hostname} ${disk.value.name}"
      size       = disk.value.size_gib * 1024 * 1024 * 1024
    }
  }

  timeouts {
    create = "20m"
  }
}

output "guest_address" {
  description = "Management IPv4 address of the MLflow guest."
  value       = local.guest_address
}
