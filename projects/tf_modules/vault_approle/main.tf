terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = ">= 5"
    }
  }
}

variable "name" {
  type = string
  description = "Name used for resources"
}

variable "member_entity_ids" {
  type = list(string)
  description = "Entity ids for the entity group"
}

variable "policies" {
  type = list(string)
  default = []
  description = "Extra policies for the approle"
}

variable "backend" {
  type = string
  description = "Backend for the approle"
}

variable "secrets" {
  type = string
  description = "Secrets backend"
}

variable "policy" {
  type = string
  description = "Policy for the approle"
}

variable "disable_yc_folder_policy" {
  type = bool
  default = false
  description = "If set, do not create a policy for Yandex Cloud"
}

locals {
    secret_path = "alwaldend.com/vault1/approles/${var.name}"
    yc_folder = "yandex.cloud/org1/folders/${replace(var.name, "_", "-")}"
}

output "secret_path" {
    description = "Secret path for this approle"
    value = local.secret_path
}

output "secret_mount" {
    description = "Secret mount"
    value = var.secrets
}

resource "vault_policy" "policy" {
  name = var.name
  policy = var.policy
}

resource "vault_policy" "yc_folder" {
  name = "${var.name}_yc_folder"
  policy = <<EOT
    path "${var.secrets}/data/${local.yc_folder}/*" {
        capabilities = ["read"]
    }
    path "${var.secrets}/metadata/${local.yc_folder}/*" {
      capabilities = ["list"]
    }
EOT
}

resource "vault_policy" "approle_secrets" {
  name = "${var.name}_approle_secrets"
  policy = <<EOT
    path "${var.secrets}/data/${local.secret_path}/*" {
       capabilities = ["create", "read", "update", "delete", "list"]
    }
    path "${var.secrets}/metadata/${local.secret_path}/*" {
      capabilities = ["list", "read", "update"]
    }
EOT
}

resource "vault_policy" "shared" {
  name = "${var.name}_shared"
  policy = <<EOT
    path "${var.secrets}/data/+/+/shared/read/approle/${var.name}/*" {
        capabilities = ["read"]
    }
    path "${var.secrets}/metadata/+/+/shared/read/approle/${var.name}/*" {
      capabilities = ["list"]
    }
EOT
}

resource "vault_approle_auth_backend_role" "role" {
  backend        = var.backend
  role_name      = var.name
  role_id        =  var.name
  secret_id_num_uses = 1
  secret_id_ttl = 3600
  token_policies = concat(
        [
            vault_policy.shared.name,
            vault_policy.approle_secrets.name,
            vault_policy.policy.name,
        ],
        var.policies,
        var.disable_yc_folder_policy ? [] : [vault_policy.yc_folder.name],
    )
}

resource "vault_policy" "policy_approle" {
  name = "${var.name}_approle"
  policy = <<EOT
    path "auth/${var.backend}/role/${vault_approle_auth_backend_role.role.role_name}/secret-id" {
        capabilities = ["update", "read"]
    }
EOT
}

resource "vault_identity_group" "group" {
  name     = var.name
  type     = "internal"
  policies = [vault_policy.policy_approle.name]
  member_entity_ids = var.member_entity_ids

  metadata = {
    version = "2"
  }
}
