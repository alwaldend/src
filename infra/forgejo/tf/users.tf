variable "vault_oauth_source_id" {
  type        = number
  description = "Verified numeric ID of the active vault OpenID Connect source in this Forgejo instance."

  validation {
    condition     = var.vault_oauth_source_id > 0 && floor(var.vault_oauth_source_id) == var.vault_oauth_source_id
    error_message = "vault_oauth_source_id must be a positive integer verified against the active vault authentication source."
  }
}

# Follow the owning login group, including its two levels of nested groups.
data "vault_identity_group" "forgejo_users" {
  group_name = "src_infra_dc1_forgejo1_users"
}

data "vault_identity_group" "forgejo_user_children" {
  for_each = toset(data.vault_identity_group.forgejo_users.member_group_ids)
  group_id = each.key
}

data "vault_identity_group" "forgejo_user_grandchildren" {
  for_each = toset(flatten([for group in data.vault_identity_group.forgejo_user_children : tolist(group.member_group_ids)]))
  group_id = each.key

  lifecycle {
    postcondition {
      condition     = length(self.member_group_ids) == 0
      error_message = "Forgejo user discovery supports two nested group levels; deeper membership must be resolved before provisioning."
    }
  }
}

locals {
  vault_user_entity_ids = toset(concat(
    tolist(data.vault_identity_group.forgejo_users.member_entity_ids),
    flatten([for group in data.vault_identity_group.forgejo_user_children : tolist(group.member_entity_ids)]),
    flatten([for group in data.vault_identity_group.forgejo_user_grandchildren : tolist(group.member_entity_ids)]),
  ))
}

data "vault_identity_entity" "forgejo_users" {
  for_each  = local.vault_user_entity_ids
  entity_id = each.key
}

locals {
  vault_user_entities = {
    for entity in data.vault_identity_entity.forgejo_users : entity.entity_name => entity
  }
}

resource "forgejo_user" "vault" {
  for_each             = local.vault_user_entities
  login                = each.key
  email                = each.value.metadata.email
  source_id            = var.vault_oauth_source_id
  login_name           = each.value.entity_id
  password             = ""
  must_change_password = false
  send_notify          = false

  lifecycle {
    prevent_destroy = true
    # OIDC owns account privilege flags. The provider invents send_notify=true
    # on import; it is create-only, so preserve imported state while disabling
    # notifications for newly created accounts above.
    ignore_changes = [admin, restricted, send_notify]

    precondition {
      condition     = each.value.metadata.username == each.key
      error_message = "The requested Forgejo username must match its authoritative Vault entity metadata."
    }
  }
}
