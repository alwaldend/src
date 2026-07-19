variable "xui_url" {
  type = string
}

variable "xui_username" {
  type = string
}

variable "xui_password" {
  type = string
}

variable "xui_token" {
  type = string
}

variable "mullvad_keys" {
  type = map(object({
    private_key = string
    address     = string
  }))
}

variable "njalla_xui_url" {
  type = string
}

variable "njalla_xui_username" {
  type = string
}

variable "njalla_xui_password" {
  type = string
}

variable "njalla_xui_token" {
  type = string
}
