terraform {
  required_providers {
    routeros = {
      source  = "terraform-routeros/routeros"
      version = "1.99.1"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.9.1"
    }
  }
}

provider "routeros" {}

provider "vault" {
  skip_child_token = true
}

resource "random_password" "printers" {
  length  = 32
  special = false
}

resource "routeros_system_user_group" "printers" {
  name    = "src_infra_printers"
  comment = "tf[infra/mikrotik/accounts]: printer network configuration"
  policy  = ["api", "read", "rest-api", "write"]
}

resource "routeros_system_user" "printers" {
  name     = "src_infra_printers"
  group    = routeros_system_user_group.printers.name
  password = random_password.printers.result
  comment  = "tf[infra/mikrotik/accounts]: printers Terraform"
  disabled = false
}

resource "vault_kv_secret_v2" "printers" {
  mount = "secrets"
  name  = "alwaldend.com/vault1/approles/src_infra_mikrotik/printers"
  data_json = jsonencode({
    mikrotik_username = routeros_system_user.printers.name
    mikrotik_password = random_password.printers.result
  })
  delete_all_versions = false
}
