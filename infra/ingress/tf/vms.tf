locals {
  cloud_init_str = file("${path.module}/../../pve/ansible/files/cloud_init.yaml")
  # cloud_init = yamldecode(local.cloud_init_str)
  vpc = {
    host1 = {
      zone = "ru-central1-d"
    }
    # host2 = {
    #   zone = "ru-central1-b"
    # }
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
  source_url = var.image_signed_url
  hardware_generation {
    generation2_features {}
  }
}

resource "yandex_vpc_network" "vpc" {
  name = "vpc"
}

resource "yandex_vpc_subnet" "vpc" {
  for_each = {
    ru-central1-d = { v4_cidr_blocks = ["10.6.0.0/24"] },
    ru-central1-b = { v4_cidr_blocks = ["10.7.0.0/24"] }
  }
  name           = "vpc-${each.key}"
  network_id     = yandex_vpc_network.vpc.id
  zone           = each.key
  v4_cidr_blocks = each.value.v4_cidr_blocks
}

resource "yandex_kms_symmetric_key" "vpc" {
  name              = "vpc"
  folder_id         = var.folder_id
  description       = "Key for disks"
  default_algorithm = "AES_128"
  rotation_period   = "8760h" // equal to 1 year
}

resource "yandex_compute_disk" "vpc" {
  for_each   = local.vpc
  name       = "vpc-${each.key}"
  type       = "network-ssd"
  zone       = each.value.zone
  kms_key_id = yandex_kms_symmetric_key.vpc.id
  size       = 10
  image_id   = yandex_compute_image.vpc.id
}

resource "yandex_compute_instance" "vpc" {
  for_each    = local.vpc
  name        = "vpc-${each.key}"
  zone        = each.value.zone
  platform_id = "standard-v3" # https://yandex.cloud/en/docs/compute/concepts/vm-platforms
  resources {
    cores         = 2
    core_fraction = 50
    memory        = 2
  }
  boot_disk {
    disk_id = yandex_compute_disk.vpc[each.key].id
  }
  network_interface {
    subnet_id      = yandex_vpc_subnet.vpc[each.value.zone].id
    nat            = true
    nat_ip_address = yandex_vpc_address.vpc[each.key].external_ipv4_address[0].address
  }
  metadata = {
    user-data = local.cloud_init_str
    # ssh-keys = "${local.cloud_init.users[0].name}:${local.cloud_init.users[0].ssh_authorized_keys[0]}}"
  }
}

output "vpc" {
  value = {
    for key, value in local.vpc : key => {
      ipv4 = yandex_vpc_address.vpc[key].external_ipv4_address[0].address
    }
  }
}
