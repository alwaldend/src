variable "runners" {
  description = "Runner allocations; each key must have an owner DNS record and Ansible inventory host."
  type = map(object({
    cpus       = optional(number, 8)
    memory_gib = optional(number, 16)
    disk_gib   = optional(number, 500)
  }))
  default = { secure = {} }

  validation {
    condition = alltrue([
      for name, runner in var.runners : can(regex("^[a-z][a-z0-9-]*$", name)) &&
      runner.cpus > 0 && floor(runner.cpus) == runner.cpus &&
      runner.memory_gib > 0 && runner.disk_gib >= 20
    ])
    error_message = "Runner names must be DNS labels with positive CPU/memory and at least 20 GiB disk."
  }
}

variable "xoa_url" {
  type        = string
  description = "Xen Orchestra WebSocket endpoint."
  default     = "wss://xoa.xcp-ng.alwaldend.com"
}

variable "xoa_template_name" {
  type        = string
  description = "Unique cloud template name within the runner resource set."
  default     = "fedora44-cloud.templates.xcp-ng.alwaldend.com"
}

variable "xoa_storage_name" {
  type        = string
  description = "Storage repository name within the runner resource set."
  default     = "Local storage"
}

variable "xoa_network_name" {
  type        = string
  description = "Network name carrying the runner subnet within its resource set."
  default     = "Pool-wide network 1"
}

variable "xoa_guest_interface" {
  type        = string
  description = "Guest interface name in the selected cloud template."
  default     = "enX0"
}

variable "gateway" {
  type        = string
  description = "Gateway for the runner subnet."
  default     = "192.168.10.1"
}

variable "dns_servers" {
  type        = list(string)
  description = "Resolvers reachable from the runner subnet."
  default     = ["192.168.10.1"]
}
