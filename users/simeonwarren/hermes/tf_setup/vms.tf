locals {
  yc_cloud_init = file("${path.module}/../../../../infra/pve/ansible/files/cloud_init.yaml")
  yc_vms = {
    host1 = {
      zone = "ru-central1-d"
    }
  }
  yc_subnets = {
    ru-central1-d = {
      v4_cidr_blocks = ["10.10.0.0/24"]
    }
  }
}

resource "yandex_dns_zone" "hermes" {
  name        = "hermes"
  description = "tf[users/simeonwarren/hermes]"
  zone        = "yc.hermes.simeonwarren.users.alwaldend.com."
  public      = true
}

resource "yandex_compute_image" "hermes" {
  name         = "hermes"
  source_image = "fd8iku26nnkveh8s4di1" # Fedora 43
}

resource "yandex_vpc_network" "hermes" {
  name = "hermes"
}

resource "yandex_vpc_subnet" "hermes" {
  for_each       = local.yc_subnets
  name           = each.key
  network_id     = yandex_vpc_network.hermes.id
  zone           = each.key
  v4_cidr_blocks = each.value.v4_cidr_blocks
}

resource "yandex_kms_symmetric_key" "hermes" {
  name              = "hermes"
  folder_id         = var.folder_id
  description       = "Key for hermes VM disks"
  default_algorithm = "AES_128"
  rotation_period   = "8760h"
}

resource "yandex_compute_disk" "hermes" {
  for_each   = local.yc_vms
  name       = each.key
  type       = "network-ssd"
  zone       = each.value.zone
  kms_key_id = yandex_kms_symmetric_key.hermes.id
  size       = 100
  image_id   = yandex_compute_image.hermes.id
}

resource "yandex_compute_instance" "hermes" {
  for_each    = local.yc_vms
  name        = each.key
  hostname    = "${each.key}.hermes.simeonwarren.users.alwaldend.com"
  zone        = each.value.zone
  platform_id = "standard-v3"
  resources {
    cores         = 2
    core_fraction = 50
    memory        = 2
  }
  boot_disk {
    disk_id = yandex_compute_disk.hermes[each.key].id
  }
  network_interface {
    subnet_id = yandex_vpc_subnet.hermes[each.value.zone].id
    nat       = true
    nat_dns_record {
      dns_zone_id = yandex_dns_zone.hermes.id
      fqdn        = "${each.key}.${yandex_dns_zone.hermes.zone}"
    }
  }
  metadata = {
    user-data          = local.yc_cloud_init
    serial-port-enable = "1"
  }
}

output "yc_vms" {
  value = {
    for key, value in local.yc_vms : key => {
      fqdn = yandex_compute_instance.hermes[key].network_interface[0].nat_dns_record[0].fqdn
      ipv4 = yandex_compute_instance.hermes[key].network_interface[0].nat_ip_address
    }
  }
}
