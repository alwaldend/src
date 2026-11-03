locals {
  dns_document = jsondecode(file("${path.module}/../dnsconfig.json"))
  name         = "local-inference"
  cloud_init   = file("${path.module}/../../cloud_init/assets/cloud_init.yaml")
  labels = {
    managed_by = "terraform"
    component  = "local-inference"
  }
}

resource "yandex_vpc_network" "host" {
  name   = local.name
  labels = local.labels
}

resource "yandex_vpc_subnet" "host" {
  name           = local.name
  zone           = var.zone
  network_id     = yandex_vpc_network.host.id
  v4_cidr_blocks = ["10.40.0.0/24"]
  labels         = local.labels
}

resource "yandex_vpc_security_group" "host" {
  name       = local.name
  network_id = yandex_vpc_network.host.id
  labels     = local.labels

  ingress {
    description    = "SSH from explicitly approved administrator networks"
    protocol       = "TCP"
    port           = 22
    v4_cidr_blocks = var.admin_ipv4_cidrs
  }

  ingress {
    description    = "Traefik HTTP challenge and HTTPS with client authentication"
    protocol       = "TCP"
    port           = 80
    v4_cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description    = "T3 Code HTTPS protected by Traefik mTLS"
    protocol       = "TCP"
    port           = 443
    v4_cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description    = "Normal host networking; Codex owns command permissions"
    protocol       = "ANY"
    v4_cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "yandex_dns_zone" "host" {
  name   = local.name
  zone   = "${local.dns_document.records.yc_ns1.NS.name}."
  public = true
  labels = local.labels
}

resource "yandex_kms_symmetric_key" "disk" {
  name                = local.name
  description         = "Encryption key for all local inference persistent storage"
  default_algorithm   = "AES_256"
  rotation_period     = "8760h"
  deletion_protection = true
  labels              = local.labels

  lifecycle {
    prevent_destroy = true
  }
}

resource "yandex_compute_disk" "boot" {
  name       = local.name
  zone       = var.zone
  type       = "network-ssd"
  size       = var.disk_size_gib
  image_id   = var.image_id
  kms_key_id = yandex_kms_symmetric_key.disk.id
  labels     = local.labels

  hardware_generation {
    legacy_features {
      pci_topology = "PCI_TOPOLOGY_V1"
    }
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "yandex_compute_instance" "host" {
  name        = local.name
  hostname    = local.name
  zone        = var.zone
  platform_id = "gpu-standard-v3"
  labels      = local.labels

  lifecycle {
    prevent_destroy = true
  }

  resources {
    cores         = 28
    core_fraction = 100
    memory        = 119
    gpus          = 1
  }

  boot_disk {
    disk_id     = yandex_compute_disk.boot.id
    auto_delete = false
  }

  network_interface {
    subnet_id          = yandex_vpc_subnet.host.id
    security_group_ids = [yandex_vpc_security_group.host.id]
    nat                = true
    nat_dns_record {
      dns_zone_id = yandex_dns_zone.host.id
      fqdn        = local.dns_document.records.t3code.CNAME.target
    }
  }

  scheduling_policy {
    preemptible = true
  }

  metadata = {
    user-data          = local.cloud_init
    serial-port-enable = "0"
  }
}

output "ansible_host" {
  description = "Set LOCAL_INFERENCE_HOST to this address for the Ansible deployment."
  value       = yandex_compute_instance.host.network_interface[0].nat_ip_address
}

output "encrypted_disk_id" {
  value = yandex_compute_disk.boot.id
}

output "kms_key_id" {
  value = yandex_kms_symmetric_key.disk.id
}
