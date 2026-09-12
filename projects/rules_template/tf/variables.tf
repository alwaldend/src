variable "dns_cloudflare_zone_id" {
  type        = string
  description = "Cloudflare zone ID supplied by the DNS Vault injector."
  default     = null
}

variable "dns_enabled" {
  type        = bool
  description = "Enable reconciliation only after this owner's documented cutover and imports."
  default     = true
}
