terraform {
  required_providers {
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "0.203.0"
    }
  }
  backend "http" {}
}

provider "yandex" { folder_id = var.folder_id }

variable "folder_id" { type = string }
variable "zone" {
  type    = string
  default = "ru-central1-d"
}
variable "image_id" {
  description = "Immutable Fedora 43 image already used by infra/threexui."
  type        = string
  default     = "fd8iku26nnkveh8s4di1"
}

locals {
  records  = jsondecode(file("${path.module}/../dnsconfig.json")).records
  hostname = trimsuffix(local.records.download_global.CNAME.target, ".")
}

resource "yandex_vpc_network" "download" { name = "download" }
resource "yandex_vpc_subnet" "download" {
  name           = "download"
  network_id     = yandex_vpc_network.download.id
  zone           = var.zone
  v4_cidr_blocks = ["10.9.0.0/24"]
}
resource "yandex_vpc_security_group" "download" {
  name       = "download"
  network_id = yandex_vpc_network.download.id
  dynamic "ingress" {
    for_each = toset([22, 80, 443])
    content {
      protocol       = "TCP"
      port           = ingress.value
      v4_cidr_blocks = ["0.0.0.0/0"]
    }
  }
  egress {
    protocol       = "ANY"
    v4_cidr_blocks = ["0.0.0.0/0"]
  }
}
resource "yandex_vpc_address" "download" {
  name = "download"
  external_ipv4_address { zone_id = var.zone }
  lifecycle { prevent_destroy = true }
}
resource "yandex_compute_disk" "content" {
  name = "download-content"
  zone = var.zone
  type = "network-hdd"
  size = 100
  lifecycle { prevent_destroy = true }
}
resource "yandex_compute_instance" "download" {
  name        = "download"
  hostname    = local.hostname
  zone        = var.zone
  platform_id = "standard-v3"
  resources {
    cores         = 2
    memory        = 2
    core_fraction = 50
  }
  boot_disk {
    auto_delete = true
    initialize_params {
      image_id = var.image_id
      size     = 20
      type     = "network-hdd"
    }
  }
  secondary_disk {
    disk_id     = yandex_compute_disk.content.id
    device_name = "download-content"
    auto_delete = false
  }
  network_interface {
    subnet_id          = yandex_vpc_subnet.download.id
    nat                = true
    nat_ip_address     = yandex_vpc_address.download.external_ipv4_address[0].address
    security_group_ids = [yandex_vpc_security_group.download.id]
  }
  metadata = {
    user-data = file("${path.module}/../../cloud_init/assets/cloud_init.yaml")
  }
}

# The allocated address remains a provider-owned fact. Global aliases use
# this delegated name, without copying the address into source or other states.
resource "yandex_dns_zone" "download" {
  name   = "download"
  zone   = "${local.records.yc_ns1.NS.name}.alwaldend.com."
  public = true
}
resource "yandex_dns_recordset" "host" {
  zone_id = yandex_dns_zone.download.id
  name    = local.records.download_global.CNAME.target
  type    = "A"
  ttl     = 300
  data    = [yandex_vpc_address.download.external_ipv4_address[0].address]
}

output "address" { value = yandex_vpc_address.download.external_ipv4_address[0].address }
output "hostname" { value = local.hostname }
