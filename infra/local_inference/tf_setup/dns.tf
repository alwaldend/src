variable "dns_cloudflare_zone_id" {
  type        = string
  default     = ""
  description = "Optional zone ID supplied by the DNS credential injector."
}

module "dns" {
  source             = "../../../projects/tf_modules/dns_records/global"
  document           = local.dns_document
  cloudflare_zone_id = var.dns_cloudflare_zone_id
  enabled            = true
}
