variable "folder_id" {
  type        = string
  description = "Yandex Cloud folder id"
}

variable "vault_url" {
  type        = string
  description = "Vault url"
  default     = "https://vault.dc1.alwaldend.com:8200"
}

variable "service_account_id" {
  type        = string
  description = "Yandex Cloud service_account_id"
}
