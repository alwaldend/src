variable "xui_url" {
  type = string
}

variable "xui_username" {
  type = string
}

variable "xui_password" {
  type = string
}

variable "mullvad_keys" {
  type = map(object({
    private_key = string
    address     = string
  }))
}

variable "node_xui_username" {
  type = string
}

variable "node_xui_base_path" {
  type = string
}

variable "node_xui_password" {
  type = string
}

variable "node_xui_tokens" {
  type = map(string)
}
