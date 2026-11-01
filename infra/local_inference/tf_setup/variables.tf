variable "folder_id" {
  type        = string
  description = "Component folder, injected through the component Vault AppRole."
}

variable "cloud_id" {
  type        = string
  description = "Cloud ID supplied by the shared Yandex authentication injector."
}

variable "zone" {
  type        = string
  default     = "ru-central1-a"
  description = "V100 availability zone; deployment still requires quota and capacity."
  validation {
    condition     = contains(["ru-central1-a", "ru-central1-b"], var.zone)
    error_message = "Broadwell V100 instances require ru-central1-a or ru-central1-b."
  }
}

variable "image_id" {
  type        = string
  description = "Immutable Ubuntu 24.04 x86_64 GPU image with a proprietary V100-compatible NVIDIA driver >= 570.26."
  validation {
    condition     = length(trimspace(var.image_id)) > 0
    error_message = "Supply a reviewed immutable GPU image ID; no moving family is selected."
  }
}

variable "admin_ipv4_cidrs" {
  type        = set(string)
  description = "Explicit IPv4 CIDRs allowed to administer the host over SSH."
  validation {
    condition = length(var.admin_ipv4_cidrs) > 0 && alltrue([
      for cidr in var.admin_ipv4_cidrs :
      can(cidrnetmask(cidr)) && try(tonumber(split("/", cidr)[1]) > 0, false)
    ])
    error_message = "Supply at least one valid IPv4 CIDR; unrestricted /0 SSH access is forbidden."
  }
}

variable "disk_size_gib" {
  type        = number
  default     = 200
  description = "Encrypted boot SSD capacity for OS, model files, and workspaces."
  validation {
    condition     = var.disk_size_gib >= 100 && floor(var.disk_size_gib) == var.disk_size_gib
    error_message = "The disk must have an integral capacity of at least 100 GiB."
  }
}
