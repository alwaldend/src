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

variable "xui_address" {
  type = string
}
