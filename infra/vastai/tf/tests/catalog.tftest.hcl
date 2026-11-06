mock_provider "vastai" {}

run "account_teams_and_shared_environment" {
  command = apply

  assert {
    condition     = length(vastai_team.managed) == 1 && vastai_team.managed["alwaldend"].name == "alwaldend"
    error_message = "The account root must preserve only the alwaldend team."
  }

  assert {
    condition     = alltrue([for name, variable in vastai_environment_variable.environment : variable.value == local.configuration.environment_variables[name]]) && keys(vastai_environment_variable.environment) == keys(local.configuration.environment_variables)
    error_message = "Account settings must read the shared environment map."
  }
}
