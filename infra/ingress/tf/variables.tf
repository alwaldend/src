variable "image_signed_url" {
  type        = string
  description = "Signed url for the image"
}

variable "wg_public_keys" {
  type        = map(string)
  description = "Wireguard public keys"
}

variable "wg_preshared_keys" {
  type        = map(string)
  description = "Wireguard preshared keys"
}

variable "wg_private_keys" {
  type        = map(string)
  description = "Wireguard interface private keys"
}

variable "folder_id" {
  type        = string
  description = "Yandex Cloud folder id"
}
