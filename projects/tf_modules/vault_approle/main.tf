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

resource "vault_policy" "policy" {
  name = var.name
  policy = var.policy
}

# Allow to read secrets on shared paths
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
  token_policies = concat([vault_policy.policy.name, vault_policy.shared.name], var.policies)
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
