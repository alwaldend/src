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
  description = "Name used for resources"
}

variable "backend" {
  type        = string
  description = "Backend for the role"
}

variable "allowed_domains" {
  type        = list(string)
  description = "Allowed domains"
}

variable "server_flag" {
  type    = bool
  default = true
}

variable "client_flag" {
  type    = bool
  default = false
}

variable "ttl" {
  type        = number
  description = "Role ttl"
  default     = 86400 * 7 # 7 days
}

variable "eab_new_member_group_ids" {
  type        = list(string)
  description = "Groups allowed to create EABs"
  default     = []
}

variable "max_ttl" {
  type        = number
  description = "Role max ttl"
  default     = 2629746 # Month
}

output "name" {
  value = vault_pki_secret_backend_role.role.name
}

output "backend" {
  value = vault_pki_secret_backend_role.role.backend
}

resource "vault_policy" "eab_new" {
  name   = "${var.backend}/${var.name}/eab_new"
  policy = <<EOT
    path "${var.backend}/roles/${var.name}/acme/new-eab" {
      capabilities = ["update"]
    }
EOT
}

resource "vault_identity_group" "eab_new" {
  name = "${var.backend}/${var.name}/eab_new"
  type = "internal"
  policies = [
    vault_policy.eab_new.name,
  ]
  member_group_ids = var.eab_new_member_group_ids
  metadata = {
    comment = "Group allowed to create EABs for ${var.backend}/roles/${var.name}"
  }
}

resource "vault_pki_secret_backend_role" "role" {
  backend            = var.backend
  name               = var.name
  ttl                = var.ttl
  max_ttl            = var.max_ttl
  allow_ip_sans      = false
  key_type           = "rsa"
  key_bits           = 4096
  allow_any_name     = false
  allow_localhost    = false
  allowed_domains    = var.allowed_domains
  allow_bare_domains = true
  allow_subdomains   = true
  server_flag        = var.server_flag
  client_flag        = var.client_flag
  no_store           = false
}

resource "vault_pki_secret_backend_role" "rsa_2048" {
  backend            = var.backend
  name               = "${var.name}_rsa_2048"
  ttl                = var.ttl
  max_ttl            = var.max_ttl
  allow_ip_sans      = false
  key_type           = "rsa"
  key_bits           = 2048
  allow_any_name     = false
  allow_localhost    = false
  allowed_domains    = var.allowed_domains
  allow_bare_domains = true
  allow_subdomains   = true
  server_flag        = var.server_flag
  client_flag        = var.client_flag
  no_store           = false
}
