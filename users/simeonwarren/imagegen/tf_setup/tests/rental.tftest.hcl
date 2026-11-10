# The actual rental module/API behavior is covered by the provider CLI E2E suite.
# This test isolates credential generation and Vault persistence in this root.
mock_provider "vault" {}
mock_provider "vastai" {}
override_module {
  target = module.imagegen
  outputs = {
    instance_id         = "71"
    https_url           = "https://203.0.113.2:27860"
    ssh_forward_command = null
  }
}

run "vault_backed_login" {
  command = apply
  variables {
    template_hash_id = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  }
  assert {
    condition     = jsondecode(vault_kv_secret_v2.web.data_json).password == random_password.web.result && jsondecode(vault_kv_secret_v2.web.data_json).username == var.web_username
    error_message = "Vault must retain the generated login used by the rental."
  }
  assert {
    condition     = length(random_password.web.result) == 48 && can(regex("^[A-Za-z0-9]+$", random_password.web.result))
    error_message = "The login must be a generated 48-character alphanumeric credential."
  }
  assert {
    condition     = vault_kv_secret_v2.web.name == "alwaldend.com/vault1/approles/user_simeonwarren/imagegen/web"
    error_message = "The login must use the regular user secret namespace."
  }
}
