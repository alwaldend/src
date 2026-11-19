provider "cloudflare" {}

variable "dns_cloudflare_zone_id" {
  type        = string
  default     = ""
  description = "Optional zone ID supplied by the DNS credential injector."
}
