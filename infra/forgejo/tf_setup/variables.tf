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
  description = "Unique cloud template name within the Forgejo resource set."
  default     = "fedora44-cloud.templates.xcp-ng.alwaldend.com"
}

variable "xoa_storage_name" {
  type        = string
  description = "Storage repository name within the Forgejo resource set."
  default     = "Local storage"
}

variable "xoa_network_name" {
  type        = string
  description = "Network name carrying the Forgejo subnet within its resource set."
  default     = "Pool-wide network 1"
}

variable "xoa_guest_interface" {
  type        = string
  description = "Guest interface name in the selected cloud template."
  default     = "enX0"
}

variable "gateway" {
  type        = string
  description = "Gateway for the Forgejo subnet."
  default     = "192.168.10.1"
}

variable "dns_servers" {
  type        = list(string)
  description = "Resolvers reachable from the Forgejo subnet."
  default     = ["192.168.10.1"]
}
