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

variable "template_name" {
  type    = string
  default = "fedora44-cloud.templates.xcp-ng.alwaldend.com"
}

variable "storage_name" {
  type    = string
  default = "Local storage"
}

variable "network_name" {
  type    = string
  default = "Pool-wide network 1"
}

data "xenorchestra_template" "download" { name_label = var.template_name }
data "xenorchestra_sr" "download" { name_label = var.storage_name }
data "xenorchestra_network" "download" { name_label = var.network_name }
data "xenorchestra_resource_set" "download" { name = "src_infra_download" }

locals {
  records  = jsondecode(file("${path.module}/../dnsconfig.json")).records
  address  = jsondecode(file("${path.module}/../../dns/dnsconfig.json")).records.download_dc1.A.address
  hostname = "${local.records.host_dc1.CNAME.name}.alwaldend.com"
  mac      = join(":", concat(["02"], [for i in range(5) : substr(sha256(local.hostname), i * 2, 2)]))
}

resource "xenorchestra_vm" "download" {
  name_label       = local.hostname
  name_description = "Static hosting managed by infra/download/local"
  template         = data.xenorchestra_template.download.id
  resource_set     = data.xenorchestra_resource_set.download.id
  cpus             = 2
  memory_max       = 2 * 1024 * 1024 * 1024
  memory_min       = 2 * 1024 * 1024 * 1024
  auto_poweron     = true
  power_state      = "Running"
  tags             = ["download", "src_infra_download"]

  cloud_config = "#cloud-config\n${yamlencode(merge(
    jsondecode(file("${path.module}/../../cloud_init/xen_linux.json")),
    { hostname = "download-local", fqdn = local.hostname },
  ))}"
  cloud_network_config = yamlencode({
    version = 1
    config = [{
      type = "physical", name = "enX0", mac_address = local.mac
      subnets = [{
        type    = "static", address = "${local.address}/24"
        gateway = "192.168.10.1", dns_nameservers = ["192.168.10.1"]
      }]
    }]
  })
  network {
    network_id  = data.xenorchestra_network.download.id
    mac_address = local.mac
  }
  disk {
    sr_id      = data.xenorchestra_sr.download.id
    name_label = "${local.hostname} boot"
    size       = 20 * 1024 * 1024 * 1024
  }
  disk {
    sr_id      = data.xenorchestra_sr.download.id
    name_label = "${local.hostname} content"
    size       = 100 * 1024 * 1024 * 1024
  }
  lifecycle {
    # This provider owns disks with the VM; never silently delete content.
    prevent_destroy = true
  }
  timeouts { create = "20m" }
}

output "address" { value = local.address }
output "hostname" { value = local.hostname }
