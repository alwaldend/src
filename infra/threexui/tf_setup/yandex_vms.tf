locals {
  yc_cloud_init = file("${path.module}/../../pve/ansible/files/cloud_init.yaml")
  yc_vms = {
    host1 = {
      zone = "ru-central1-d"
    }
  }
  yc_subnets = {
    ru-central1-d = {
      v4_cidr_blocks = ["10.8.0.0/24"]
    }
  }
}

resource "yandex_dns_zone" "threexui" {
  name        = "threexui"
  description = "tf[infra/threexui]"
  zone        = "yc.threexui.alwaldend.com."
  public      = true
}

resource "yandex_compute_image" "threexui" {
  name         = "threexui"
  source_image = "fd8iku26nnkveh8s4di1" # Fedora 43
}

resource "yandex_vpc_network" "threexui" {
  name = "threexui"
}

resource "yandex_vpc_subnet" "threexui" {
  for_each       = local.yc_subnets
  name           = each.key
  network_id     = yandex_vpc_network.threexui.id
  zone           = each.key
  v4_cidr_blocks = each.value.v4_cidr_blocks
}

resource "yandex_kms_symmetric_key" "threexui" {
  name              = "threexui"
  folder_id         = var.folder_id
  description       = "Key for threexui VM disks"
  default_algorithm = "AES_128"
  rotation_period   = "8760h"
}

resource "yandex_compute_disk" "threexui" {
  for_each   = local.yc_vms
  name       = each.key
  type       = "network-ssd"
  zone       = each.value.zone
  kms_key_id = yandex_kms_symmetric_key.threexui.id
  size       = 20
  image_id   = yandex_compute_image.threexui.id
}

resource "yandex_compute_instance" "threexui" {
  for_each    = local.yc_vms
  name        = each.key
  hostname    = "${each.key}.nodes.threexui.alwaldend.com"
  zone        = each.value.zone
  platform_id = "standard-v3"
  resources {
    cores         = 2
    core_fraction = 50
    memory        = 2
  }
  boot_disk {
    disk_id = yandex_compute_disk.threexui[each.key].id
  }
  network_interface {
    subnet_id = yandex_vpc_subnet.threexui[each.value.zone].id
    nat       = true
    nat_dns_record {
      dns_zone_id = yandex_dns_zone.threexui.id
      fqdn        = "${each.key}.nodes.${yandex_dns_zone.threexui.zone}"
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
      fqdn = yandex_compute_instance.threexui[key].network_interface[0].nat_dns_record[0].fqdn
      ipv4 = yandex_compute_instance.threexui[key].network_interface[0].nat_ip_address
    }
  }
}
