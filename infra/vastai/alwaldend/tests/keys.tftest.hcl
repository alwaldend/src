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
}
