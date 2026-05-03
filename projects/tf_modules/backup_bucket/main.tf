terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = ">= 5"
    }
    yandex = {
      source  = "yandex-cloud/yandex"
      version = ">= 0.201.0"
    }
  }
}

variable "name" {
  type = string
  description = "Name used for resources"
}

variable "service_accounts" {
    type = list(string)
    description = "List of service account ids to grant access to"
}

variable "folder_id" {
  type = string
  description = "Folder id"
}

variable "bucket_prefix" {
  type = string
  description = "Bucket prefix"
}

variable "secret_mount" {
  type = string
  default = "secrets"
  description = "Secret mount"
}

variable "secret_name" {
  type = string
  description = "Secret name"
}

locals {
    labels = {
        "name": var.name,
        "managed_by": "terraform",
        "module": "projects_tf_modules_backup_bucket"
    }
}

resource "yandex_kms_symmetric_key" "key" {
  name              = var.name
  folder_id = var.folder_id
  labels = local.labels
  description       = "Key for ${var.name}"
  default_algorithm = "AES_128"
  rotation_period   = "8760h" // equal to 1 year
}

resource "yandex_storage_bucket" "bucket" {
  bucket = "${var.bucket_prefix}-${var.name}"
  folder_id = var.folder_id
  tags = local.labels
  max_size = 1024 * 1024 * 1024 * 1000
  default_storage_class = "ICE"
  versioning {
    enabled = true
  }
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = yandex_kms_symmetric_key.key.id
        sse_algorithm     = "aws:kms"
      }
    }
  }
}

resource "yandex_storage_bucket_grant" "grant" {
  bucket = yandex_storage_bucket.bucket.bucket
  dynamic "grant" {
    for_each = {for id in var.service_accounts : id => id}
    content {
        id = grant.value
        permissions = ["READ", "WRITE"]
        type        = "CanonicalUser"
    }
 }
}

resource "vault_kv_secret_v2" "secret" {
  mount               = var.secret_mount
  name                = var.secret_name
  data_json = jsonencode(
    {
        id = yandex_storage_bucket.bucket.id
        folder_id = var.folder_id
        bucket  = yandex_storage_bucket.bucket.bucket
    }
  )
}
