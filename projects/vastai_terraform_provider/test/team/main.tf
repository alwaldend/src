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

resource "vastai_ssh_key" "test_ssh_key" {
  public_key = "${trimspace(file("${path.module}/test_public_key.pub"))} test_vastai_terraform_provider_${var.test_marker}"
}

resource "vastai_environment_variable" "test_environment_variable" {
  name  = "TEST_VASTAI_TERRAFORM_PROVIDER"
  value = var.test_marker
}

resource "vastai_team_role" "test_team_role" {
  name        = "test_vastai_terraform_provider"
  permissions = jsonencode({ api = { user_read = {} } })
}
