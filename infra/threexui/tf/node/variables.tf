variable "xui_url" {
  type = string
}

variable "xui_username" {
  type = string
}

variable "xui_password" {
  type = string
}

variable "xui_base_path" {
  type = string
}

variable "mullvad_key" {
  type = object({
    private_key = string
    address     = string
  })
}

variable "mullvad_min" {
  type        = list(string)
  description = "List of name of wireguard servers"
}

variable "http_proxies" {
  type        = map(string)
  description = "HTTP proxy accounts (username -> password)"
}

variable "mikrotik_parent_proxy" {
  type = object({
    port     = number
    outbound = string
  })
  default     = null
  description = "Optional unauthenticated HTTP parent proxy dedicated to MikroTik"
}

variable "xui_address" {
  type = string
}
