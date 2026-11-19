variable "dns_cloudflare_zone_id" {
  type        = string
  description = "Cloudflare zone containing this project's public DNS records."
  default     = null
}

variable "dns_enabled" {
  type        = bool
  description = "Enable DNS management after the authorized ownership cutover and import."
  default     = true
}
