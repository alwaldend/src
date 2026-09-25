mock_provider "yandex" {
  mock_resource "yandex_vpc_address" {
    defaults = {
      external_ipv4_address = [{ address = "192.0.2.10" }]
    }
  }
}

variables { folder_id = "fixture-folder" }

# Exercise compute, retained storage, firewall and allocated-address DNS as
# one graph. The mock provider cannot contact Yandex or create resources.
run "cloud_provisioning_graph" {
  command = plan
}
