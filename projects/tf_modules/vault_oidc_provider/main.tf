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
  description = "Name"
}

variable "issuer_host" {
  type        = string
  description = "Issuer host"
}

variable "scopes_supported" {
  type        = list(string)
  description = "List of OIDC scopes"
  default     = null
}

variable "group_ids" {
  type        = list(string)
  description = "Group ids for the assignment"
  default     = null
}

variable "allowed_read_clients_group_ids" {
  type        = list(string)
  description = "Group ids allowed to read the client"
  default     = null
}

variable "redirect_urls" {
  type        = list(string)
  description = "Redirect urls for the client"
  default     = null
}

variable "verification_ttl" {
  type        = number
  description = "Verification ttl for the key"
  default     = 3600 * 2 # Two hours
}

variable "rotation_period" {
  type        = number
  description = "Rotation period for the key"
  default     = 3600 # Hour
}

resource "vault_identity_oidc_assignment" "assignment" {
  name      = var.name
  group_ids = var.group_ids
}

resource "vault_identity_oidc_key" "key" {
  name               = var.name
  algorithm          = "RS512"
  verification_ttl   = var.verification_ttl
  rotation_period    = var.rotation_period
  allowed_client_ids = ["*"]
}

resource "vault_identity_oidc_client" "client" {
  name          = var.name
  redirect_uris = var.redirect_urls
  key           = vault_identity_oidc_key.key.name
  assignments = [
    vault_identity_oidc_assignment.assignment.name,
  ]
  id_token_ttl     = vault_identity_oidc_key.key.rotation_period
  access_token_ttl = vault_identity_oidc_key.key.verification_ttl
}

resource "vault_policy" "allowed_read_client" {
  name   = "${var.name}_allowed_read_client"
  policy = <<EOT
    path "identity/oidc/client/${vault_identity_oidc_client.client.name}" {
      capabilities = [ "read" ]
    }
EOT
}

resource "vault_identity_group" "allowed_read_client" {
  name = "${var.name}_allowed_read_client"
  type = "internal"
  policies = [
    vault_policy.allowed_read_client.name,
  ]
  member_group_ids = var.allowed_read_clients_group_ids
  metadata = {
    comment = "Allowed to read client info of ${vault_identity_oidc_client.client.name}"
  }
}

resource "vault_identity_oidc_provider" "provider" {
  name          = var.name
  https_enabled = true
  issuer_host   = var.issuer_host
  allowed_client_ids = [
    vault_identity_oidc_client.client.client_id,
  ]
  scopes_supported = var.scopes_supported
}
