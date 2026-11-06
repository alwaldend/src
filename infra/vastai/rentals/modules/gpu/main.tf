terraform {
  required_providers {
    vastai = {
      source  = "alwaldend/vastai"
      version = "0.0.0-dev"
    }
  }
}

variable "image" {
  description = "Container image for the rental."
  type        = string
  default     = "ubuntu:24.04"
  nullable    = false
}

variable "disk_gb" {
  description = "Disk allocation included when comparing total hourly offer prices."
  type        = number
  default     = 1
  nullable    = false
  validation {
    condition     = var.disk_gb > 0
    error_message = "Disk allocation must be positive."
  }
}

variable "maximum_hourly_price" {
  description = "Maximum accepted total hourly offer price in USD."
  type        = number
  default     = 1
  nullable    = false
  validation {
    condition     = var.maximum_hourly_price > 0
    error_message = "The hourly price ceiling must be positive."
  }
}

variable "label" {
  description = "Account-visible label for the rented contract."
  type        = string
  default     = "infra_vastai_rental"
  nullable    = false
}

variable "release" {
  description = "Skip marketplace discovery when releasing the existing contract."
  type        = bool
  default     = false
}

data "vastai_offers" "cheapest" {
  count = var.release ? 0 : 1
  query = jsonencode({
    rentable          = { eq = true }
    rented            = { eq = false }
    num_gpus          = { gte = 1 }
    disk_space        = { gte = var.disk_gb }
    allocated_storage = var.disk_gb
    type              = "on-demand"
    order             = [["dph_total", "asc"]]
    limit             = 1
  })
  lifecycle {
    postcondition {
      condition     = length(self.offers) > 0
      error_message = "No rentable on-demand GPU offers are available."
    }
  }
}

resource "vastai_instance" "worker" {
  count    = var.release ? 0 : 1
  offer_id = data.vastai_offers.cheapest[0].offers[0].id
  image    = var.image
  disk_gb  = var.disk_gb
  label    = var.label

  lifecycle {
    ignore_changes = [offer_id]
    precondition {
      condition     = data.vastai_offers.cheapest[0].offers[0].hourly_price <= var.maximum_hourly_price
      error_message = "The cheapest available offer exceeds the configured hourly price ceiling."
    }
  }
}

output "selected_offer" {
  description = "Current cheapest matching marketplace offer; an existing rental retains its original offer."
  value       = try(data.vastai_offers.cheapest[0].offers[0], null)
}

output "rental_contract_id" {
  value = try(vastai_instance.worker[0].id, null)
}

output "rental_hourly_price" {
  value = try(data.vastai_offers.cheapest[0].offers[0].hourly_price, null)
}

output "ssh" {
  value = try({
    host = vastai_instance.worker[0].ssh_host
    port = vastai_instance.worker[0].ssh_port
  }, null)
}
