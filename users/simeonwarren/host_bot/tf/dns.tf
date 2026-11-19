variable "dns_enabled" {
  type        = bool
  description = "Enable DNS management after importing and reviewing this owner's existing records."
  default     = true
}

variable "dns_cloudflare_zone_id" {
  type        = string
  description = "Cloudflare zone identifier supplied by the DNS credential injector."
  default     = ""
}

variable "dns_routeros_hosturl" {
  type        = string
  description = "RouterOS management endpoint supplied by the DNS credential injector."
  default     = ""
}

variable "dns_routeros_username" {
  type        = string
  description = "RouterOS DNS username supplied by the DNS credential injector."
  sensitive   = true
  default     = ""
}

variable "dns_routeros_password" {
  type        = string
  description = "RouterOS DNS password supplied by the DNS credential injector."
  sensitive   = true
  default     = ""
}

provider "cloudflare" {
}

provider "routeros" {
  alias    = "dns"
  hosturl  = var.dns_routeros_hosturl
  username = var.dns_routeros_username
  password = var.dns_routeros_password
}

module "dns" {
  source = "../../../../projects/tf_modules/dns_records"

  providers = {
    cloudflare = cloudflare
    routeros   = routeros.dns
  }

  document           = jsondecode(file("${path.module}/../dnsconfig.json"))
  cloudflare_zone_id = var.dns_cloudflare_zone_id
  enabled            = var.dns_enabled
}
