terraform {
  required_version = ">= 1.0"
  backend "http" {}
  required_providers {
    vastai = {
      source  = "alwaldend/vastai"
      version = "0.0.0-dev"
    }
  }
}

provider "vastai" {}

variable "test_marker" {
  description = "Non-secret test marker used to exercise in-place updates."
  type        = string
  default     = "initial"
}

variable "test_team_api_key" {
  description = "Generated team credential supplied by the project's Vault injector."
  type        = string
  sensitive   = true
}

resource "vastai_api_key" "test_api_key" {
  name         = "test_vastai_terraform_provider_reader"
  team_api_key = var.test_team_api_key
  permissions  = jsonencode({ api = { user_read = {} } })
}

resource "vastai_environment_variable" "test_environment_variable" {
  name  = "TEST_VASTAI_TERRAFORM_PROVIDER"
  value = var.test_marker
}

resource "vastai_team_role" "test_team_role" {
  name        = "test_vastai_terraform_provider"
  permissions = jsonencode({ api = { "*" = { "api.users.current" = { GET = {} } } } })
}
