# Failure case: provisioning the standalone identity must not require the
# parent root's context, state outputs, or a not-yet-created DNS policy.
mock_provider "vault" {
  mock_data "vault_auth_backend" {
    defaults = { accessor = "auth_approle_fixture" }
  }
  mock_data "vault_identity_entity" {
    defaults = { entity_id = "operator-fixture-id" }
  }
}

run "standalone_identity" {
  command = plan
}
