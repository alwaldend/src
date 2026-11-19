# Lookups run with the deployment's own identity, which has access only to
# its resource set. The provider rejects ambiguous name matches.
data "xenorchestra_template" "openhands" {
  name_label = var.xoa_template_name
}

data "xenorchestra_sr" "openhands" {
  name_label = var.xoa_storage_name
}

data "xenorchestra_network" "openhands" {
  name_label = var.xoa_network_name
}

data "xenorchestra_resource_set" "openhands" {
  name = "src_infra_openhands"
}
