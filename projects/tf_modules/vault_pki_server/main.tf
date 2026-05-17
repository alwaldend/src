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

variable "ttl" {
  type        = number
  description = "Role ttl"
  default     = 86400 * 7 # 7 days
}

variable "max_ttl" {
  type        = number
  description = "Role max ttl"
  default     = 2629746 # Month
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
  server_flag        = true
  client_flag        = false
  no_store           = false
}
