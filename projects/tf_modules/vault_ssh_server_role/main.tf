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
  description = "Role backend"
}

variable "allow_subdomains" {
  type = bool
  description = "If set, allow subdomains"
  default = true
}

variable "allowed_domains" {
  type        = string
  description = "Allowed domains"
}

variable "ttl" {
  type        = number
  description = "Role ttl"
  default     = 31536000 // Year
}

output "policy" {
  description = "Policy name"
  value       = vault_policy.policy.name
}

resource "vault_ssh_secret_backend_role" "role" {
  backend                 = var.backend
  name                    = var.name
  max_ttl                 = var.ttl
  ttl                     = var.ttl
  key_type                = "ca"
  allow_host_certificates = true
  allow_subdomains        = var.allow_subdomains
  allow_bare_domains      = true
  allowed_domains         = var.allowed_domains
}

resource "vault_policy" "policy" {
  name   = var.name
  policy = <<EOT
    path "${var.backend}/sign/${vault_ssh_secret_backend_role.role.name}" {
      capabilities = ["update"]
    }
EOT
}
