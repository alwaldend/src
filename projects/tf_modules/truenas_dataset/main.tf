terraform {
  required_providers {
    truenas = {
      source  = "truenas/truenas"
      version = "1.9.3"
    }
  }
}

variable "parent" {
  description = "Existing parent dataset name, or the name output of another dataset module."
  type        = string

  validation {
    condition     = alltrue([for component in split("/", var.parent) : !contains(["", ".", ".."], component)])
    error_message = "The parent must be a dataset path without empty or dot components."
  }
}

variable "config" {
  description = "Child name, account ownership and optional filesystem/share overrides."
  type = object({
    name        = string
    owner       = optional(string, "root")
    traverse    = optional(bool, false)
    access      = optional(map(string), {})
    compression = optional(string, "lz4")
    atime       = optional(string, "OFF")
    copies      = optional(number)
    quota       = optional(number)
    smb = optional(object({
      name  = optional(string)
      hosts = optional(list(string), [])
    }))
  })

  validation {
    condition     = !contains(["", ".", ".."], var.config.name) && !strcontains(var.config.name, "/")
    error_message = "The child name must be one non-empty dataset component."
  }

  validation {
    condition     = alltrue([for access in values(var.config.access) : contains(["MODIFY", "TRAVERSE"], access)])
    error_message = "Named access must be MODIFY or TRAVERSE."
  }
}

resource "truenas_dataset" "this" {
  name               = "${var.parent}/${var.config.name}"
  type               = "filesystem"
  compression        = var.config.compression
  atime              = var.config.atime
  acltype            = "nfsv4"
  aclmode            = "PASSTHROUGH"
  copies             = var.config.copies
  quota              = var.config.quota
  inherit_encryption = true
  force_destroy      = false

  lifecycle {
    prevent_destroy = true
  }
}

data "truenas_user" "owner" {
  username = var.config.owner
}

data "truenas_user" "access" {
  for_each = var.config.access
  username = each.key
}

locals {
  modify_access = { for name, access in var.config.access : name => access if access == "MODIFY" }
}

resource "truenas_filesystem_acl" "this" {
  path = truenas_dataset.this.mountpoint
  uid  = data.truenas_user.owner.uid
  entries = jsonencode(concat(
    [
      { tag = "owner@", type = "ALLOW", perms = { BASIC = "FULL_CONTROL" }, flags = { BASIC = "INHERIT" } },
      { tag = "USER", id = data.truenas_user.owner.uid, type = "ALLOW", perms = { BASIC = "FULL_CONTROL" }, flags = { BASIC = "INHERIT" } },
    ],
    var.config.traverse ? [{ tag = "everyone@", type = "ALLOW", perms = { BASIC = "TRAVERSE" }, flags = { BASIC = "NOINHERIT" } }] : [],
    [for name, access in var.config.access : {
      tag   = "USER", id = data.truenas_user.access[name].uid, type = "ALLOW"
      perms = { BASIC = access }
      flags = { BASIC = access == "MODIFY" ? "INHERIT" : "NOINHERIT" }
    }],
  ))
  recursive = length(local.modify_access) > 0
  traverse  = false
}

resource "truenas_smb_share" "this" {
  count     = var.config.smb == null ? 0 : 1
  name      = coalesce(var.config.smb.name, var.config.name)
  path      = truenas_dataset.this.mountpoint
  purpose   = "DEFAULT_SHARE"
  browsable = false
  options = {
    hostsallow = var.config.smb.hosts
    hostsdeny  = length(var.config.smb.hosts) == 0 ? [] : ["ALL"]
  }

  depends_on = [truenas_filesystem_acl.this]
}

resource "truenas_smb_share_acl" "this" {
  count      = var.config.smb == null ? 0 : 1
  share_name = truenas_smb_share.this[0].name
  share_acl = concat(
    [{ ae_who_id = { id_type = "USER", id = data.truenas_user.owner.uid }, ae_type = "ALLOWED", ae_perm = "FULL" }],
    [for name in keys(local.modify_access) : {
      ae_who_id = { id_type = "USER", id = data.truenas_user.access[name].uid }
      ae_type   = "ALLOWED"
      ae_perm   = "CHANGE"
    }],
  )
}

output "name" { value = truenas_dataset.this.name }
output "mountpoint" { value = truenas_dataset.this.mountpoint }
