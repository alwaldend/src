terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = ">= 5"
    }
  }
}

variable "name" {
  type        = string
  description = "Name for resources"
}

variable "backend" {
  type        = string
  description = "Transit backend path"
}


variable "users_member_group_ids" {
  type        = list(string)
  description = "Member groups ids allowed to encrypt data"
  default     = []
}

variable "decryptors_member_group_ids" {
  type        = list(string)
  description = "Member groups ids allowed to decrypt data"
  default     = []
}

resource "vault_transit_secret_backend_key" "key" {
  backend = var.backend
  name    = var.name
}

resource "vault_identity_group" "users" {
  name = "${var.name}_users"
  type = "internal"
  policies = [
    resource.vault_policy.user.name,
  ]
  member_group_ids = concat(var.users_member_group_ids, var.decryptors_member_group_ids)
  metadata = {
    comment = "Allowed to encrypt data"
  }
}

resource "vault_policy" "user" {
  name   = "${var.name}_user"
  policy = <<EOT
    path "${var.backend}/encrypt/${vault_transit_secret_backend_key.key.name}" {
      capabilities = ["update"]
    }
EOT
}

resource "vault_identity_group" "decryptors" {
  name = "${var.name}_decryptors"
  type = "internal"
  policies = [
    resource.vault_policy.decryptor.name,
  ]
  member_group_ids = var.decryptors_member_group_ids
  metadata = {
    comment = "Allowed to encrypt data"
  }
}

resource "vault_policy" "decryptor" {
  name   = "${var.name}_decryptor"
  policy = <<EOT
    path "${var.backend}/decrypt/${vault_transit_secret_backend_key.key.name}" {
      capabilities = ["update"]
    }
EOT
}
