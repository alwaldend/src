mock_provider "vastai" {}
mock_provider "vault" {
  mock_data "vault_generic_secret" {
    defaults = { data_json = "{\"metadata\":{\"sshpubkey\":\"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEZpeHR1cmVGaXh0dXJlRml4dHVyZUZpeHR1cmVG fixture\"}}" }
  }
}

run "account_teams_and_shared_environment" {
  command = apply
  assert {
    condition = alltrue([for name, key in vastai_ssh_key.host :
      key.public_key == jsondecode(data.vault_generic_secret.host_ssh[name].data_json).metadata.sshpubkey
    ]) && length(vastai_ssh_key.host) > 0
    error_message = "Account SSH registrations must use Vault user metadata."
  }


  assert {
    condition     = length(vastai_team.managed) == 1 && vastai_team.managed["alwaldend"].name == "alwaldend"
    error_message = "The account root must preserve only the alwaldend team."
  }

  assert {
    condition     = alltrue([for name, variable in vastai_environment_variable.environment : variable.value == local.configuration.environment_variables[name]]) && keys(vastai_environment_variable.environment) == keys(local.configuration.environment_variables)
    error_message = "Account settings must read the shared environment map."
  }
}
