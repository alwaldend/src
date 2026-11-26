variable "flux_url" {
  type        = string
  description = "Flux url"
  default     = "https://operator.flux.alwaldend.com"
}

variable "vault_host" {
  type        = string
  description = "Vault host for the Flux OIDC issuers"
  default     = "vault.alwaldend.com:8200"
}
