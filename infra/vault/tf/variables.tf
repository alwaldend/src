variable "folder_id" {
  type        = string
  description = "Yandex Cloud folder id"
}

variable "vault_url" {
  type        = string
  description = "Vault url"
  default     = "https://vault.alwaldend.com:8200"
}

variable "pve1_url" {
  type        = string
  description = "Pve url"
  default     = "https://pve.alwaldend.com:8006"
}

variable "xoa_url" {
  type        = string
  description = "Xen Orchestra browser URL"
  default     = "https://xoa.xcp-ng.alwaldend.com"
}

variable "forgejo_url" {
  type        = string
  description = "Forgejo url"
  default     = "https://git.alwaldend.com"
}

variable "harbor_url" {
  type        = string
  description = "Harbor url"
  default     = "https://harbor.alwaldend.com"
}

variable "flux_url" {
  type        = string
  description = "Flux url"
  default     = "https://operator.flux.alwaldend.com"
}

variable "vault_host" {
  type        = string
  description = "Vault host"
  default     = "vault.alwaldend.com:8200"
}

variable "service_account_id" {
  type        = string
  description = "Yandex Cloud service_account_id"
}
