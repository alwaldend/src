# The standalone AppRole root must be applied before these membership lookups.
data "vault_identity_entity" "src_infra_download" {
  entity_name = "src_infra_download"
}

data "vault_identity_group" "src_infra_download" {
  group_name = "src_infra_download"
}
