terraform {
  backend "http" {}
  required_providers {
    vastai = {
      source  = "alwaldend/vastai"
      version = "0.0.0-dev"
    }
  }
}

provider "vastai" {}

locals {
  disk_gb = 3.5
}

data "vastai_offers" "cheapest" {
  query = jsonencode({
    rentable          = { eq = true }
    rented            = { eq = false }
    num_gpus          = { gte = 1 }
    disk_space        = { gte = local.disk_gb }
    allocated_storage = local.disk_gb
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

resource "vastai_instance" "test" {
  offer_id = data.vastai_offers.cheapest.offers[0].id
  image    = "ubuntu:24.04"
  disk_gb  = local.disk_gb
  label    = "infra_vastai_rentals_test"

  lifecycle {
    ignore_changes = [offer_id]

    precondition {
      condition     = data.vastai_offers.cheapest.offers[0].hourly_price <= 1
      error_message = "The cheapest available offer exceeds USD 1 per hour."
    }
  }
}

output "selected_offer" {
  description = "Current cheapest marketplace offer; an existing contract retains its original offer."
  value       = data.vastai_offers.cheapest.offers[0]
}

output "rental_contract_id" {
  value = vastai_instance.test.id
}
