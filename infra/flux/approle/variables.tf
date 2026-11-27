variable "vault_host" {
  type        = string
  description = "Vault host for the Flux OIDC issuers"
  default     = "vault.alwaldend.com:8200"
}
