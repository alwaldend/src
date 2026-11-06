terraform {
  required_providers {
    vastai = {
      source = "alwaldend/vastai"
    }
  }
}

variable "offer_id" {
  type        = number
  default     = null
  description = "Optional reviewed on-demand offer ID; otherwise select the cheapest matching datacenter offer."

  validation {
    condition     = var.offer_id == null ? true : var.offer_id > 0 && floor(var.offer_id) == var.offer_id
    error_message = "offer_id must be a positive integer."
  }
}

variable "min_gpu_ram_gb" {
  type        = number
  default     = 24
  description = "Minimum VRAM per GPU for automatic offer selection, in GiB."
  nullable    = false

  validation {
    condition     = var.min_gpu_ram_gb >= 8
    error_message = "min_gpu_ram_gb must be at least 8."
  }
}

variable "max_hourly_price" {
  type        = number
  default     = 1
  description = "Maximum API-reported hourly USD price for automatic selection, including the requested disk allocation. This is not a total spending limit."
  nullable    = false

  validation {
    condition     = var.max_hourly_price > 0
    error_message = "max_hourly_price must be positive."
  }
}

variable "gpu_name" {
  type        = string
  default     = null
  description = "Optional exact GPU model filter for automatic selection, for example RTX 4090."
}

data "vastai_offers" "imagegen" {
  query = jsonencode(merge({
    datacenter        = { eq = true }
    verified          = { eq = true }
    rentable          = { eq = true }
    rented            = { eq = false }
    num_gpus          = { eq = 1 }
    gpu_ram           = { gte = var.min_gpu_ram_gb * 1024 }
    reliability       = { gte = 0.99 }
    direct_port_count = { gte = 1 }
    disk_space        = { gte = var.disk_gb }
    allocated_storage = var.disk_gb
    dph_total         = { lte = var.max_hourly_price }
    type              = "on-demand"
    order             = [["dph_total", "asc"]]
    limit             = 1
  }, var.gpu_name == null ? {} : { gpu_name = { eq = var.gpu_name } }, var.offer_id == null ? {} : { id = { eq = var.offer_id } }))

  lifecycle {
    postcondition {
      condition     = length(self.offers) > 0
      error_message = "No verified datacenter GPU offer matches the requested memory, storage, connectivity, and hourly price limit."
    }
    postcondition {
      condition     = alltrue([for offer in self.offers : offer.hourly_price <= var.max_hourly_price && offer.gpu_ram_mb >= var.min_gpu_ram_gb * 1024 && offer.num_gpus == 1])
      error_message = "The returned offer exceeds the hourly price limit or does not meet GPU requirements."
    }
    postcondition {
      condition     = var.offer_id == null ? true : alltrue([for offer in self.offers : offer.id == var.offer_id])
      error_message = "The returned offer does not match the requested offer ID."
    }
  }
}

variable "template_hash_id" {
  type        = string
  description = "Reviewed AI-Dock Forge template hash with all published port options removed."
  nullable    = false

  validation {
    condition     = can(regex("^[0-9a-f]{32}$", var.template_hash_id))
    error_message = "template_hash_id must be a 32-character lowercase hexadecimal hash."
  }
}

variable "disk_gb" {
  type        = number
  default     = 100
  description = "Local storage in GB. Vast recommends at least 50 GB, or 100 GB for multiple models."
  nullable    = false

  validation {
    condition     = var.disk_gb > 0
    error_message = "disk_gb must be positive."
  }
}

variable "label" {
  type        = string
  default     = "imagegen"
  description = "Mutable instance label."
  nullable    = false
}

variable "state" {
  type        = string
  default     = "running"
  description = "Desired power state. Stopped instances still incur storage charges."
  nullable    = false

  validation {
    condition     = contains(["running", "stopped"], var.state)
    error_message = "state must be running or stopped."
  }
}

variable "local_port" {
  type        = number
  default     = 7860
  description = "Local loopback port used by the SSH forwarding command; no remote port is published."
  nullable    = false

  validation {
    condition     = var.local_port >= 1024 && var.local_port <= 65535 && floor(var.local_port) == var.local_port
    error_message = "local_port must be an integer from 1024 to 65535."
  }
}

data "vastai_template" "imagegen" {
  hash_id = var.template_hash_id
}

resource "vastai_instance" "imagegen" {
  offer_id         = data.vastai_offers.imagegen.offers[0].id
  template_hash_id = var.template_hash_id
  image            = data.vastai_template.imagegen.image
  disk_gb          = var.disk_gb
  label            = var.label
  state            = var.state
  runtype          = "ssh_direct"

  # AI-Dock init.sh accepts KEY=value arguments; ___ becomes a literal space.
  # Vast replaces the image entrypoint in SSH mode, so explicitly start Forge.
  onstart = join(" ", [
    "/opt/ai-dock/bin/init.sh",
    "CF_QUICK_TUNNELS=false",
    "SUPERVISOR_NO_AUTOSTART=caddy,cloudflared,jupyter,quicktunnel,serviceportal,sshd,syncthing",
    "AUTO_UPDATE=false",
    "FORGE_PORT_LOCAL=17860",
    "FORGE_ARGS=--server-name___127.0.0.1",
  ])

  lifecycle {
    # Marketplace changes must not replace an existing paid contract.
    ignore_changes = [offer_id]

    precondition {
      condition     = data.vastai_template.imagegen.environment_only
      error_message = "Template Docker options must contain only environment assignments. Remove all -p, -P, --publish, networking, and other Docker flags from a reviewed copy of the Forge template."
    }
  }
}

output "instance_id" {
  value       = vastai_instance.imagegen.id
  description = "Rental contract ID."
}

output "ssh_host" {
  value       = vastai_instance.imagegen.ssh_host
  description = "SSH host; null until observed by a refresh."
}

output "ssh_port" {
  value       = vastai_instance.imagegen.ssh_port
  description = "Mapped SSH port; null until observed by a refresh."
}

output "ssh_forward_command" {
  value = (vastai_instance.imagegen.ssh_host == null || vastai_instance.imagegen.ssh_port == null) ? null : format(
    "ssh -N -o ExitOnForwardFailure=yes -L 127.0.0.1:%d:127.0.0.1:17860 -p %d root@'%s'",
    var.local_port, vastai_instance.imagegen.ssh_port, replace(vastai_instance.imagegen.ssh_host, "'", "'\"'\"'"),
  )
  description = "Run locally after SSH is ready, then open local_url. Uses the caller's SSH agent or normal SSH identity configuration."
}

output "local_url" {
  value       = "http://127.0.0.1:${var.local_port}"
  description = "Forge URL available only while the SSH forwarding command is running."
}
