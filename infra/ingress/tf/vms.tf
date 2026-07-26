locals {
  vpc = {
    host1 = {
      zone = "ru-central1-d"
    }
    host2 = {
      zone = "ru-central1-b"
    }
  }
}

resource "yandex_vpc_address" "vpc" {
  for_each = local.vpc
  name     = "vpc-${each.key}"
  external_ipv4_address {
    zone_id = each.value.zone
  }
}

resource "yandex_compute_image" "vpc" {
  name       = "vpc"
  source_url = "https://download.fedoraproject.org/pub/fedora/linux/releases/44/Cloud/x86_64/images/Fedora-Cloud-Base-Generic-44-1.7.x86_64.qcow2"
}

resource "yandex_vpc_network" "vpc" {
  name = "vpc"
}

resource "yandex_vpc_subnet" "vpc" {
  network_id     = yandex_vpc_network.vpc.id
  zone           = local.vpc.host1.zone
  v4_cidr_blocks = ["10.5.0.0/24"]
}

resource "yandex_compute_disk" "vpc" {
  for_each = local.vpc
  name     = "vpc-${each.key}"
  type     = "network-ssd"
  size     = 20
  zone     = each.value.zone
  image_id = yandex_compute_image.vpc.id
}

resource "yandex_compute_instance" "vpc" {
  for_each = local.vpc
  name     = "vpc-${each.key}"
  zone     = each.value.zone
  resources {
    cores  = 1
    memory = 2
  }
  boot_disk {
    disk_id = yandex_compute_disk.vpc[each.key].id
  }
  network_interface {
    subnet_id      = yandex_vpc_subnet.vpc.id
    nat            = true
    nat_ip_address = yandex_vpc_address.vpc[each.key].external_ipv4_address[0].address
  }
  metadata = {
    user-data = file("${path.module}/../../pve/ansible/files/cloud_init.yaml")
  }
}

output "vpc" {
  value = {
    for key, value in local.vpc : key => {
      ipv4 = yandex_vpc_address.vpc[key].external_ipv4_address[0].address
    }
  }
}
