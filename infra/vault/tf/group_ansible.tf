resource "vault_identity_group" "ansible" {
  name = "ansible"
  type = "internal"
  policies = [
    vault_policy.ssh_clients_sign_admins.name,
  ]
  member_group_ids = [
    data.vault_identity_group.standalone_approles["src_infra_dev"].group_id,
    data.vault_identity_group.src_infra_download.group_id,
    module.src_infra_dc1_forgejo1_approle.group_id,
    module.src_infra_local_inference_approle.group_id,
    module.src_infra_harbor_approle.group_id,
    module.src_infra_flux_approle.group_id,
    module.src_infra_threexui_approle.group_id,
    module.src_infra_ingress_approle.group_id,
    module.src_infra_openhands.group_id,
    data.vault_identity_group.standalone_approles["src_infra_grafana"].group_id,
    data.vault_identity_group.standalone_approles["src_infra_mlflow"].group_id,
    data.vault_identity_group.standalone_approles["src_infra_opencanary"].group_id,
    data.vault_identity_group.standalone_approles["src_infra_victoria_metrics"].group_id,
    module.users_simeonwarren_approle.group_id,
  ]
  metadata = {
    comment = "Group with access to ansible"
  }
}
