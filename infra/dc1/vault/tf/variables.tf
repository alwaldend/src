variable "folder_id" {
  type        = string
  description = "Yandex Cloud folder id"
}

variable "vault_url" {
  type        = string
  description = "Vault url"
  default     = "https://vault.dc1.alwaldend.com:8200"
}

variable "pve1_url" {
  type        = string
  description = "Pve url"
  default     = "https://host1.pve1.dc1.alwaldend.com:8006"
}

variable "forgejo_url" {
  type        = string
  description = "Forgejo url"
  default     = "https://forgejo.alwaldend.com"
}

variable "harbor_url" {
  type        = string
  description = "Harbor url"
  default     = "https://harbor.alwaldend.com"
}

variable "vault_host" {
  type        = string
  description = "Vault host"
  default     = "vault.dc1.alwaldend.com:8200"
}

variable "service_account_id" {
  type        = string
  description = "Yandex Cloud service_account_id"
}
