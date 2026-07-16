resource "vault_identity_group" "approles" {
  name = "approles"
  type = "internal"
  policies = [
    vault_policy.tf_token.name,
    vault_policy.pki_server.name,
  ]
  member_entity_ids = [
    module.src_infra_dc1_consul1_approle.entity_id,
    module.src_infra_dc1_forgejo1_approle.entity_id,
    module.src_infra_dc1_pve1_approle.entity_id,
    module.src_infra_dc1_router1_approle.entity_id,
    module.src_infra_dc1_vault_approle.entity_id,
    module.src_infra_dns_approle.entity_id,
    module.src_infra_yandex_cloud_org1_approle.entity_id,
    module.src_infra_harbor_approle.entity_id,
    module.src_infra_flux_approle.entity_id,
    module.src_infra_flux_git_approle.entity_id,
    module.src_infra_flux_cluster_approle.entity_id,
    module.src_infra_threexui_approle.entity_id,
    module.src_third_party_approle.entity_id,
    module.src_projects_alwaldend_com_approle.entity_id,
    module.users_simeonwarren_approle.entity_id,
  ]
  metadata = {
    comment = "All approles"
  }
}
