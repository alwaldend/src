# Lookups run with the deployment's own identity, which has access only to
# its resource set. The provider rejects ambiguous name matches.
data "xenorchestra_template" "forgejo" {
  name_label = var.xoa_template_name
}

data "xenorchestra_sr" "forgejo" {
  name_label = var.xoa_storage_name
}

data "xenorchestra_network" "forgejo" {
  name_label = var.xoa_network_name
}
