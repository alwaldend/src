variable "name" {
  type        = string
  description = "Name used for the AppRole, entity, group, SSH role, and PKI role"
}

variable "member_entity_ids" {
  type        = list(string)
  description = "Entity ids allowed to assume the AppRole"
  default     = []
}

variable "secrets" {
  type        = string
  description = "Secrets backend"
}

variable "approle_backend" {
  type        = string
  description = "AppRole auth backend"
}

variable "approle_backend_accessor" {
  type        = string
  description = "AppRole auth backend accessor"
}

variable "ssh_backend" {
  type        = string
  description = "SSH secrets backend"
}

variable "ssh_allowed_domains" {
  type        = string
  description = "Domain the AppRole may sign host keys for"
}

variable "pki_backend" {
  type        = string
  description = "Server PKI backend"
}

variable "pki_allowed_domains" {
  type        = list(string)
  description = "Domains the PKI server role may issue certificates for"
}
