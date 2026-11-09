mock_provider "vastai" {}
mock_provider "vault" {}

run "team_keys_and_environment" {
  command = apply

  assert {
    condition     = alltrue([for name, variable in vastai_environment_variable.environment : variable.value == local.configuration.environment_variables[name]]) && keys(vastai_environment_variable.environment) == keys(local.configuration.environment_variables)
    error_message = "The team must read all settings from the shared configuration."
  }

  assert {
    condition     = vastai_api_key.keys["rentals"].permissions == ""
    error_message = "Generated team keys must request default full access."
  }

  assert {
    condition     = vault_kv_secret_v2.keys["rentals"].name == "alwaldend.com/vault1/approles/src_infra_vastai/vastai/alwaldend/rentals"
    error_message = "Rental access must be stored separately from the operator-owned root team key."
  }

  assert {
    condition     = vastai_api_key.imagegen.name == "users_simeonwarren_imagegen" && vastai_api_key.imagegen.permissions == ""
    error_message = "The existing alwaldend root must own the dedicated image-generation key without changing its name."
  }

  assert {
    condition     = vault_kv_secret_v2.imagegen.name == "alwaldend.com/vault1/approles/user_simeonwarren/imagegen/vastai" && jsondecode(vault_kv_secret_v2.imagegen.data_json).api_key == vastai_api_key.imagegen.key
    error_message = "The team root must store the generated key in the regular user namespace."
  }
}
