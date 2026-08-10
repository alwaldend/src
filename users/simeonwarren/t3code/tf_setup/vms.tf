locals {
  yc_cloud_init = file("${path.module}/../../../../infra/pve/ansible/files/cloud_init.yaml")
  yc_vms = {
    host1 = {
      zone = "ru-central1-d"
    }
  }
  yc_subnets = {
    ru-central1-d = {
      v4_cidr_blocks = ["10.30.0.0/24"]
    }
  }
}

resource "yandex_dns_zone" "t3code" {
  name        = "t3code"
  description = "tf[users/simeonwarren/t3code]"
  zone        = "yc.t3code.simeonwarren.users.alwaldend.com."
  public      = true
}

resource "yandex_compute_image" "t3code" {
  name         = "t3code"
  source_image = "fd8iku26nnkveh8s4di1" # Fedora 43
}

resource "yandex_vpc_network" "t3code" {
  name = "t3code"
}

resource "yandex_vpc_subnet" "t3code" {
  for_each       = local.yc_subnets
  name           = each.key
  network_id     = yandex_vpc_network.t3code.id
  zone           = each.key
  v4_cidr_blocks = each.value.v4_cidr_blocks
}

resource "yandex_kms_symmetric_key" "t3code" {
  name              = "t3code"
  folder_id         = var.folder_id
  description       = "Key for t3code VM disks"
  default_algorithm = "AES_128"
  rotation_period   = "8760h"
}

resource "yandex_compute_disk" "t3code" {
  for_each   = local.yc_vms
  name       = each.key
  type       = "network-ssd"
  zone       = each.value.zone
  kms_key_id = yandex_kms_symmetric_key.t3code.id
  size       = 100
  image_id   = yandex_compute_image.t3code.id
}

resource "yandex_compute_instance" "t3code" {
  for_each    = local.yc_vms
  name        = each.key
  hostname    = "${each.key}.t3code.simeonwarren.users.alwaldend.com"
  zone        = each.value.zone
  platform_id = "standard-v3"
  resources {
    cores         = 2
    core_fraction = 50
    memory        = 2
  }
  boot_disk {
    disk_id = yandex_compute_disk.t3code[each.key].id
  }
  network_interface {
    subnet_id = yandex_vpc_subnet.t3code[each.value.zone].id
    nat       = true
    nat_dns_record {
      dns_zone_id = yandex_dns_zone.t3code.id
      fqdn        = "${each.key}.${yandex_dns_zone.t3code.zone}"
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
      fqdn = yandex_compute_instance.t3code[key].network_interface[0].nat_dns_record[0].fqdn
      ipv4 = yandex_compute_instance.t3code[key].network_interface[0].nat_ip_address
    }
  }
}
