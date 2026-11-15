terraform {
  backend "http" {}
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = ">= 5"
    }
    vastai = {
      source  = "alwaldend/vastai"
      version = "0.0.0-dev"
    }
  }
}

provider "vastai" {}

provider "vault" {
  skip_child_token = true
}

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

# Read the canonical user identity and select its public key for registration.
data "vault_generic_secret" "host_ssh" {
  for_each = toset(local.configuration.ssh_key_entities)
  path     = "identity/entity/name/${each.key}"
}

resource "vastai_ssh_key" "host" {
  for_each   = data.vault_generic_secret.host_ssh
  public_key = try(jsondecode(each.value.data_json).metadata.sshpubkey, "")

  lifecycle {
    precondition {
      condition     = can(regex("^(ssh-ed25519|ssh-rsa|ecdsa-sha2-nistp(256|384|521)) [A-Za-z0-9+/]+={0,2}( [^\\r\\n]*)?$", try(jsondecode(each.value.data_json).metadata.sshpubkey, "")))
      error_message = "The configured Vault user metadata must contain one OpenSSH public key in sshpubkey."
    }
  }
}
