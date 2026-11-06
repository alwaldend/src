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

variable "component_vault_token" {
  description = "Component AppRole token injected by AL for its own Vault destinations."
  type        = string
  sensitive   = true
}

provider "vault" {
  alias            = "component"
  token            = var.component_vault_token
  skip_child_token = true
}

resource "vastai_environment_variable" "enable_https" {
  name  = "ENABLE_HTTPS"
  value = "true"
}

locals {
  teams = jsondecode(file("${path.module}/teams.json"))
  secret_targets = merge({}, [
    for team_name, team in local.teams : {
      for key_name, key in try(team.api_keys, {}) : "${team_name}/${key_name}" => {
        name = key.secret_name
        team = team_name
        key  = key_name
        auth = try(key.vault_auth, "provider_project")
      }
    }
  ]...)
  valid_secret_targets = length(distinct([for target in values(local.secret_targets) : target.name])) == length(local.secret_targets) && alltrue([
    for target in values(local.secret_targets) : contains(["default", "provider_project"], target.auth)
  ])
}

resource "vastai_team" "teams" {
  for_each = { for name, team in local.teams : name => team if !can(team.api_keys) }
  name     = try(each.value.name, each.key)

  lifecycle {
    prevent_destroy = true
  }
}

module "teams" {
  for_each = { for name, team in local.teams : name => team if can(team.api_keys) }
  source   = "../../../projects/tf_modules/vastai_team"
  name     = try(each.value.name, each.key)
  api_keys = {
    for name, key in each.value.api_keys : name => {
      name        = key.name
      permissions = key.permissions == null ? "" : jsonencode(key.permissions)
    }
  }
}

resource "vault_kv_secret_v2" "project_tokens" {
  for_each = { for name, target in local.secret_targets : name => target if target.auth != "default" }
  mount    = "secrets"
  name     = each.value.name
  data_json = jsonencode({
    api_key = module.teams[each.value.team].api_keys[each.value.key]
  })
  delete_all_versions = false

  lifecycle {
    precondition {
      condition     = local.valid_secret_targets
      error_message = "Each key must select a known Vault identity and a distinct secret destination."
    }
  }
}

resource "vault_kv_secret_v2" "component_tokens" {
  provider = vault.component
  for_each = { for name, target in local.secret_targets : name => target if target.auth == "default" }
  mount    = "secrets"
  name     = each.value.name
  data_json = jsonencode({
    api_key = module.teams[each.value.team].api_keys[each.value.key]
  })
  delete_all_versions = false

  lifecycle {
    precondition {
      condition     = local.valid_secret_targets
      error_message = "Each key must select a known Vault identity and a distinct secret destination."
    }
  }
}
