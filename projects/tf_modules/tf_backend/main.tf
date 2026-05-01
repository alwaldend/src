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

resource "yandex_kms_symmetric_key" "key" {
  name              = var.name
  folder_id = var.folder_id
  description       = "Key for ${var.name}"
  default_algorithm = "AES_128"
  rotation_period   = "8760h" // equal to 1 year
}

resource "yandex_storage_bucket" "bucket" {
  bucket = "${var.bucket_prefix}-${var.name}"
  folder_id = var.folder_id
  max_size = 1024 * 1024 * 1024
  default_storage_class = "COLD"
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

resource "yandex_iam_service_account" "account" {
  name        = var.name
  folder_id = var.folder_id
  description = "Terraform backend for ${var.name}"
}

resource "yandex_resourcemanager_folder_iam_member" "member" {
  role      = "editor"
  folder_id = var.folder_id
  member    = "serviceAccount:${yandex_iam_service_account.account.id}"
}

resource "yandex_iam_service_account_static_access_key" "key" {
  service_account_id = yandex_iam_service_account.account.id
  description        = "Static key for ${var.name}"
}

resource "yandex_storage_bucket_grant" "grant" {
  bucket = yandex_storage_bucket.bucket.bucket
  grant {
    id =  yandex_iam_service_account.account.id
    permissions = ["READ", "WRITE"]
    type        = "CanonicalUser"
 }
}

resource "vault_kv_secret_v2" "secret" {
  mount               = var.secret_mount
  name                = var.secret_name
  data_json = jsonencode(
    {
        name  = yandex_storage_bucket.bucket.bucket
        id = yandex_storage_bucket.bucket.id
        tags = yandex_storage_bucket.bucket.tags
        access_key = yandex_iam_service_account_static_access_key.key.access_key
        secret_key = yandex_iam_service_account_static_access_key.key.secret_key

    }
  )
}
