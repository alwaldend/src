terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "5.22.0"
    }
  }
  backend "http" {}
}

provider "cloudflare" {}

variable "dns_cloudflare_zone_id" {
  type        = string
  default     = ""
  description = "Optional zone ID supplied by the DNS credential injector."
}

module "dns" {
  source             = "../../../projects/tf_modules/dns_records/global"
  document           = jsondecode(file("${path.module}/../dnsconfig.json"))
  cloudflare_zone_id = var.dns_cloudflare_zone_id
  enabled            = true
}
