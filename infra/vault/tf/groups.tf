resource "vault_identity_group" "sre" {
  name     = "sre"
  type     = "internal"
  policies = []
  member_entity_ids = [
    vault_identity_entity.simeonwarren.id,
  ]
  member_group_ids = [
    vault_identity_group.dev.id,
    vault_identity_group.ansible.id,
  ]
  metadata = {
    comment = "Generic admin group"
  }
}

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
