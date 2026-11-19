variable "dns_enabled" {
  description = "Enable only with the coordinated ownership transfer revision."
  type        = bool
  default     = true
}

variable "dns_cloudflare_zone_id" {
  type    = string
  default = ""
}

variable "dns_routeros_hosturl" {
  type    = string
  default = ""
}

variable "dns_routeros_username" {
  type      = string
  default   = ""
  sensitive = true
}

variable "dns_routeros_password" {
  type      = string
  default   = ""
  sensitive = true
}

module "dns" {
  source             = "../../../projects/tf_modules/dns_records"
  document           = jsondecode(file("${path.module}/../dnsconfig.json"))
  cloudflare_zone_id = var.dns_cloudflare_zone_id
  enabled            = var.dns_enabled
  providers = {
    cloudflare = cloudflare
    routeros   = routeros.dns
  }
}

output "dns_import_addresses" {
  value = module.dns.import_addresses
}
