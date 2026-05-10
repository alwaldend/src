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
  description = "Common name"
}

variable "backend" {
  type        = string
  description = "Signing backend"
}

variable "backend_cert" {
  type        = string
  description = "Backend certificate"
}

variable "mount_path" {
  type        = string
  description = "mount path"
}

variable "default_lease" {
  type    = number
  default = 3600 # Hour
}

variable "max_lease" {
  type    = number
  default = 86400 # Day
}

variable "ca_ttl" {
  type        = number
  description = "TTL for the ca certificate"
  default     = 31536000 # Year
}

output "backend" {
  description = "Generated backend path"
  value       = vault_mount.mount.path
}

output "certificate" {
  description = "CA certificate"
  value       = vault_pki_secret_backend_intermediate_set_signed.signed.certificate
}

resource "vault_mount" "mount" {
  path                      = var.mount_path
  type                      = "pki"
  description               = "${var.name}, signed internally"
  default_lease_ttl_seconds = var.default_lease
  max_lease_ttl_seconds     = var.max_lease
}

resource "vault_pki_secret_backend_intermediate_cert_request" "request" {
  depends_on  = [vault_mount.mount]
  backend     = vault_mount.mount.path
  type        = "internal"
  common_name = var.name
  key_type    = "rsa"
  key_bits    = "4096"
}

resource "vault_pki_secret_backend_root_sign_intermediate" "sign" {
  depends_on           = [vault_pki_secret_backend_intermediate_cert_request.request]
  backend              = var.backend
  csr                  = vault_pki_secret_backend_intermediate_cert_request.request.csr
  common_name          = var.name
  exclude_cn_from_sans = true
  max_path_length      = 1
  ttl                  = var.ca_ttl
}

resource "vault_pki_secret_backend_intermediate_set_signed" "signed" {
  depends_on = [vault_pki_secret_backend_root_sign_intermediate.sign]
  backend    = vault_mount.mount.path
  certificate = format(
    "%s\n%s",
    vault_pki_secret_backend_root_sign_intermediate.sign.certificate,
    var.backend_cert,
  )
}

resource "vault_pki_secret_backend_config_issuers" "issuer" {
  depends_on                    = [vault_pki_secret_backend_intermediate_set_signed.signed]
  backend                       = vault_mount.mount.path
  default                       = vault_pki_secret_backend_intermediate_set_signed.signed.imported_issuers[0]
  default_follows_latest_issuer = true
}
