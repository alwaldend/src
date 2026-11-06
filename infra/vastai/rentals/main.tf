terraform {
  backend "http" {}
  required_providers {
    vastai = {
      source  = "alwaldend/vastai"
      version = "0.0.0-dev"
    }
  }
}

variable "alwaldend_api_key" {
  description = "Dedicated alwaldend team key injected from Vault."
  type        = string
  sensitive   = true
}

provider "vastai" {
  alias   = "alwaldend"
  api_key = var.alwaldend_api_key
}

variable "image" {
  description = "Container image override; null uses the module default."
  type        = string
  default     = null
}

variable "disk_gb" {
  description = "Disk allocation override; null uses the module default."
  type        = number
  default     = null
}

variable "maximum_hourly_price" {
  description = "Total hourly price ceiling override; null uses the module default."
  type        = number
  default     = null
}

variable "label" {
  description = "Contract label override; null uses the module default."
  type        = string
  default     = null
}

variable "release" {
  description = "Release the existing contract without marketplace discovery."
  type        = bool
  default     = false
}

module "gpu" {
  source    = "./modules/gpu"
  providers = { vastai = vastai.alwaldend }

  image                = var.image
  disk_gb              = var.disk_gb
  maximum_hourly_price = var.maximum_hourly_price
  label                = var.label
  release              = var.release
}

output "selected_offer" {
  value = module.gpu.selected_offer
}

output "rental_contract_id" {
  value = module.gpu.rental_contract_id
}

output "rental_hourly_price" {
  value = module.gpu.rental_hourly_price
}

output "ssh" {
  value = module.gpu.ssh
}
