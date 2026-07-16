resource "vault_identity_group" "allowed_entity_lookup" {
  name = "allowed_entity_lookup"
  type = "internal"
  policies = [
    vault_policy.allowed_entity_lookup.name,
  ]
  member_group_ids = [
    module.src_infra_dc1_pve1_approle.group_id,
    module.src_infra_yandex_cloud_org1_approle.group_id,
  ]
  metadata = {
    comment = "Allowed to read entity info, required for data.vault_identity_group (https://registry.terraform.io/providers/hashicorp/vault/latest/docs/data-sources/identity_group)"
  }
}

resource "vault_policy" "allowed_entity_lookup" {
  name   = "allowed_entity_lookup"
  policy = <<EOT
    path "identity/lookup/entity" {
      capabilities = ["update"]
    }
EOT
}
