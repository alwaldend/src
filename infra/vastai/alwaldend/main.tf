terraform {
  backend "http" {}
  required_providers {
    vastai = {
      source  = "alwaldend/vastai"
      version = "0.0.0-dev"
    }
    vault = {
      source  = "hashicorp/vault"
      version = ">= 5"
    }
  }
}

provider "vastai" {}

provider "vault" {
  skip_child_token = true
}

locals {
  configuration = jsondecode(file("${path.module}/../config.json"))
  api_keys      = toset(local.configuration.teams.alwaldend.api_keys)
}

resource "vastai_environment_variable" "environment" {
  for_each = local.configuration.environment_variables
  name     = each.key
  value    = each.value
}

resource "vastai_api_key" "keys" {
  for_each    = local.api_keys
  name        = "infra_vastai_alwaldend_${each.key}"
  permissions = ""

  lifecycle {
    precondition {
      condition     = can(regex("^[a-z][a-z0-9_]*$", each.key))
      error_message = "Invalid team key name: use lowercase letters, digits and underscores."
    }
  }
}

resource "vault_kv_secret_v2" "keys" {
  for_each = vastai_api_key.keys
  mount    = "secrets"
  name     = "alwaldend.com/vault1/approles/src_infra_vastai/vastai/alwaldend/${each.key}"
  data_json = jsonencode({
    api_key = each.value.key
  })
  delete_all_versions = false
}

resource "vastai_api_key" "imagegen" {
  name        = local.configuration.teams.alwaldend.imagegen_api_key_name
  permissions = ""
}

resource "vault_kv_secret_v2" "imagegen" {
  mount = "secrets"
  name  = "alwaldend.com/vault1/approles/user_simeonwarren/imagegen/vastai"
  data_json = jsonencode({
    api_key = vastai_api_key.imagegen.key
  })
  delete_all_versions = false
}
