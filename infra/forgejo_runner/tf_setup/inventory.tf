data "xenorchestra_resource_set" "runner" {
  name = "src_infra_forgejo_runner"
}

# Lookups run with the deployment's own identity, which has access only to
# its resource set. The provider rejects ambiguous name matches.
data "xenorchestra_template" "runner" {
  name_label = var.xoa_template_name
}

data "xenorchestra_sr" "runner" {
  name_label = var.xoa_storage_name
}

data "xenorchestra_network" "runner" {
  name_label = var.xoa_network_name
}
