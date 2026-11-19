variable "xoa_url" {
  type        = string
  description = "Xen Orchestra WebSocket endpoint."
  default     = "wss://xoa.xcp-ng.alwaldend.com"
}

variable "xoa_insecure" {
  type        = bool
  description = "Explicit bootstrap override for an XO certificate that is not yet trusted."
  default     = false
}

variable "xoa_template_name" {
  type        = string
  description = "Unique cloud template name within the OpenHands resource set."
  default     = "fedora44-cloud.templates.xcp-ng.alwaldend.com"
}

variable "xoa_storage_name" {
  type        = string
  description = "Storage repository name within the OpenHands resource set."
  default     = "Local storage"
}

variable "xoa_network_name" {
  type        = string
  description = "Network name carrying the OpenHands subnet within its resource set."
  default     = "Pool-wide network 1"
}

variable "xoa_guest_interface" {
  type        = string
  description = "Guest interface name in the selected cloud template."
  default     = "enX0"
}

variable "gateway" {
  type        = string
  description = "Gateway for the OpenHands subnet."
  default     = "192.168.10.1"
}

variable "dns_servers" {
  type        = list(string)
  description = "Resolvers reachable from the OpenHands subnet."
  default     = ["192.168.10.1"]
}

variable "vms" {
  description = "OpenHands components to create, keyed by the dnsconfig record prefix owning the hostname and address."
  type = map(object({
    dns_key      = string
    description  = string
    tags         = list(string)
    cpus         = number
    memory_bytes = number
    disk_bytes   = number
  }))
  default = {
    canvas = {
      dns_key      = "host1_canvas"
      description  = "OpenHands Agent Canvas managed by infra/openhands/tf_setup"
      tags         = ["openhands", "canvas", "src_infra_openhands"]
      cpus         = 2
      memory_bytes = 4 * 1024 * 1024 * 1024
      disk_bytes   = 20 * 1024 * 1024 * 1024
    }
    server = {
      dns_key      = "host1_server"
      description  = "OpenHands agent server managed by infra/openhands/tf_setup"
      tags         = ["openhands", "server", "src_infra_openhands"]
      cpus         = 4
      memory_bytes = 8 * 1024 * 1024 * 1024
      disk_bytes   = 60 * 1024 * 1024 * 1024
    }
    automation = {
      dns_key      = "host1_automation"
      description  = "OpenHands automation server managed by infra/openhands/tf_setup"
      tags         = ["openhands", "automation", "src_infra_openhands"]
      cpus         = 2
      memory_bytes = 4 * 1024 * 1024 * 1024
      disk_bytes   = 20 * 1024 * 1024 * 1024
    }
  }
}
