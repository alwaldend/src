variable "folder_id" {
  type        = string
  description = "Yandex Cloud folder id"
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
