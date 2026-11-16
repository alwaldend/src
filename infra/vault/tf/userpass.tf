locals {
  userpass_users = toset([
    vault_identity_entity_alias.userpass_simeonwarrenbot.name,
  ])
  userpass_paths = {
    for username in local.userpass_users : username => "alwaldend.com/users/${username}/userpass"
  }
}

resource "random_password" "userpass" {
  for_each = local.userpass_users
  length   = 72
  special  = false
}

resource "vault_kv_secret_v2" "userpass" {
  for_each = local.userpass_paths
  mount    = vault_mount.secrets.path
  name     = each.value
  data_json = jsonencode({
    username = each.key
    password = random_password.userpass[each.key].result
  })
  delete_all_versions = false
}

resource "vault_policy" "userpass_writer" {
  name   = "userpass_writer"
  policy = <<EOT
%{for path in values(local.userpass_paths)~}
    path "${vault_mount.secrets.path}/data/${path}" {
      capabilities = ["create", "read", "update", "delete"]
    }
    path "${vault_mount.secrets.path}/metadata/${path}" {
      capabilities = ["read", "update"]
    }
%{endfor~}
EOT
}

# Preserve unreadable password fields while detecting account deletion.
resource "vault_generic_endpoint" "userpass" {
  for_each             = local.userpass_users
  path                 = "auth/${vault_auth_backend.userpass.path}/users/${each.key}"
  ignore_absent_fields = true
  data_json = jsonencode({
    password = jsondecode(vault_kv_secret_v2.userpass[each.key].data_json).password
  })
}
