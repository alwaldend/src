terraform {
  backend "http" {}
  required_providers {
    vastai = {
      source  = "alwaldend/vastai"
      version = "0.0.0-dev"
    }
  }
}

provider "vastai" {}

locals {
  configuration = jsondecode(file("${path.module}/../config.json"))
}

resource "vastai_environment_variable" "environment" {
  for_each = local.configuration.environment_variables
  name     = each.key
  value    = each.value
}

resource "vastai_team" "managed" {
  for_each = local.configuration.teams
  name     = each.key

  lifecycle {
    prevent_destroy = true
  }
}
