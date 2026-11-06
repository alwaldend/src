mock_provider "vastai" {}

mock_provider "vault" {
  mock_resource "vault_kv_secret_v2" {
    defaults = { id = "fixture_project_secret" }
  }
}

mock_provider "vault" {
  alias = "component"
  mock_resource "vault_kv_secret_v2" {
    defaults = { id = "fixture_component_secret" }
  }
}

variables {
  component_vault_token = "fixture_component_token"
}

run "catalog_team_keys_and_vault_identities" {
  command = apply

  assert {
    condition     = contains(keys(module.teams["alwaldend"].api_keys), "rentals")
    error_message = "The team catalog must create the rental key through its team module."
  }

  assert {
    condition     = vault_kv_secret_v2.component_tokens["alwaldend/rentals"].id == "fixture_component_secret"
    error_message = "The rental secret must use the component Vault identity."
  }

  assert {
    condition     = vault_kv_secret_v2.component_tokens["alwaldend/rentals"].name == "alwaldend.com/vault1/approles/src_infra_vastai/vastai/alwaldend"
    error_message = "The rental secret must remain at the component's configured destination."
  }

  assert {
    condition     = vault_kv_secret_v2.project_tokens["test_provider/test_project"].id == "fixture_project_secret"
    error_message = "The existing project secret must retain its address and Vault identity."
  }
}
