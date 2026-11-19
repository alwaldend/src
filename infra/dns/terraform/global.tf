terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "5.22.0"
    }
  }
}

module "dns" {
  source             = "../../../projects/tf_modules/dns_records/global"
  document           = jsondecode(file("${path.module}/dnsconfig.json"))
  cloudflare_zone_id = var.dns_cloudflare_zone_id
  enabled            = var.dns_enabled
}
