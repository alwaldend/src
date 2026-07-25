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

# resource "yandex_compute_instance" "vpc" {
#   for_each = local.vpc
#   name        = "vpc-${each.key}"
#   platform_id = "standard-v3"
#   zone        = "ru-central1-a"
#
#   resources {
#     cores  = 2
#     memory = 4
#   }
#
#   boot_disk {
#     disk_id = yandex_compute_disk.boot-disk.id
#   }
#
#   network_interface {
#     index     = 1
#     subnet_id = yandex_vpc_subnet.foo.id
#   }
#
#   metadata = {
#     foo      = "bar"
#     ssh-keys = "ubuntu:${file("~/.ssh/id_ed25519.pub")}"
#   }
# }

output "vpc" {
  value = {
    for key, value in local.vpc : key => {
      ipv4 = yandex_vpc_address.vpc[key].external_ipv4_address[0].address
    }
  }
}
