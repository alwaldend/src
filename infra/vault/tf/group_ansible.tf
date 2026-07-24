resource "vault_identity_group" "ansible" {
  name = "ansible"
  type = "internal"
  policies = [
    vault_policy.ssh_clients_sign_admins.name,
  ]
  member_group_ids = [
    module.src_infra_dc1_forgejo1_approle.group_id,
    module.src_infra_harbor_approle.group_id,
    module.src_infra_flux_approle.group_id,
    module.src_infra_threexui_approle.group_id,
    module.users_simeonwarren_approle.group_id,
  ]
  metadata = {
    comment = "Group with access to ansible"
  }
}
