variable "image_signed_url" {
  type        = string
  description = "Signed url for the image"
}

variable "wg_public_key" {
  type        = string
  description = "Wireguard public key"
}

variable "wg_preshared_key" {
  type        = string
  description = "Wireguard preshared key"
}

variable "wg_interface_private_key" {
  type        = string
  description = "Wireguard interface private key"
}

variable "folder_id" {
  type        = string
  description = "Yandex Cloud folder id"
}
