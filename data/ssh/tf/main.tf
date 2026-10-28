terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "5.8.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "2.9.0"
    }
  }
  backend "local" {}
}

# Refresh reads public keys using the caller's existing Vault credentials.
provider "vault" {
  skip_child_token = true
}

variable "output_directory" {
  type        = string
  description = "Absolute path to the data/ssh directory to refresh."
  nullable    = false

  validation {
    condition     = startswith(var.output_directory, "/")
    error_message = "output_directory must be an absolute path."
  }
}

data "vault_generic_secret" "ca" {
  for_each              = toset(["clients", "servers"])
  path                  = "ssh/${each.key}/config/ca"
  with_lease_start_time = false
}

resource "local_file" "ca" {
  for_each        = data.vault_generic_secret.ca
  filename        = "${var.output_directory}/ca_${each.key}.pub"
  content         = nonsensitive(each.value.data["public_key"])
  file_permission = "0644"
}
