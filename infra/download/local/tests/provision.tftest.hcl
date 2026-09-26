mock_provider "xenorchestra" {
  mock_data "xenorchestra_template" {
    defaults = { id = "00000000-0000-4000-8000-000000000001" }
  }
  mock_data "xenorchestra_sr" {
    defaults = { id = "00000000-0000-4000-8000-000000000002" }
  }
  mock_data "xenorchestra_network" {
    defaults = { id = "00000000-0000-4000-8000-000000000003" }
  }
  mock_data "xenorchestra_resource_set" {
    defaults = { id = "00000000-0000-4000-8000-000000000004" }
  }
}

# Exercise cloud-init, canonical DNS inputs, disks and the provider schema as
# one graph, without XO authentication or any live provider operation.
run "local_provisioning_graph" {
  command = plan
}
