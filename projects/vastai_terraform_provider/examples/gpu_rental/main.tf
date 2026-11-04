terraform {
  required_providers {
    vastai = { source = "alwaldend/vastai" }
  }
}

provider "vastai" {}

variable "offer_id" {
  description = "A currently available on-demand offer ID selected after reviewing the offer search."
  type        = number
}

data "vastai_offers" "gpu" {
  query = jsonencode({
    verified = { eq = true }
    rentable = { eq = true }
    rented   = { eq = false }
    num_gpus = { eq = 1 }
    type     = "on-demand"
    order    = [["dph_total", "asc"]]
    limit    = 5
  })
}

resource "vastai_instance" "worker" {
  offer_id = var.offer_id
  image    = "ubuntu:24.04"
  disk_gb  = 20
  label    = "terraform-worker"
}

output "offers" {
  value = data.vastai_offers.gpu.offers
}

output "ssh" {
  value = {
    host = vastai_instance.worker.ssh_host
    port = vastai_instance.worker.ssh_port
  }
}
