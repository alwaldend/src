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

variable "cloud_id" {
  type = string
  description = "Cloud id"
}

variable "secret_name" {
  type = string
  description = "Base path for secrets"
}

variable "secret_mount" {
  type = string
  default = "secrets"
  description = "Mount for secrets"
}

resource "yandex_resourcemanager_folder" "folder" {
    name = var.name
    cloud_id = var.cloud_id
    description = "Folder for ${var.name} (Managed by terraform)"
}

resource "yandex_iam_service_account" "account" {
  name        = var.name
  folder_id = yandex_resourcemanager_folder.folder.id
  description = "Editor account for the folder ${var.name}"
}

resource "yandex_resourcemanager_folder_iam_member" "member" {
  role      = "editor"
  folder_id = yandex_resourcemanager_folder.folder.id
  member    = "serviceAccount:${yandex_iam_service_account.account.id}"
}

resource "yandex_iam_service_account_static_access_key" "static_key" {
  service_account_id = yandex_iam_service_account.account.id
  description        = "Static key for ${var.name}"
}

resource "yandex_iam_service_account_key" "iam_key" {
  service_account_id = yandex_iam_service_account.account.id
  description        = "Iam key for ${var.name}"
  key_algorithm      = "RSA_4096"
}

resource "vault_kv_secret_v2" "folder" {
  mount               = var.secret_mount
  name                = "${var.secret_name}/folder"
  data_json = jsonencode(
    {
      id = yandex_resourcemanager_folder.folder.id
      name = yandex_resourcemanager_folder.folder.name
      cloud_id = yandex_resourcemanager_folder.folder.cloud_id
      folder_id = yandex_resourcemanager_folder.folder.folder_id
    }
  )
}

resource "vault_kv_secret_v2" "account" {
  mount               = var.secret_mount
  name                = "${var.secret_name}/account"
  data_json = jsonencode(
    {
      id = yandex_iam_service_account.account.id
      name = yandex_iam_service_account.account.name
      cloud_id = yandex_resourcemanager_folder.folder.cloud_id
      folder_id = yandex_resourcemanager_folder.folder.folder_id
      service_account_id = yandex_iam_service_account.account.service_account_id
    }
  )
}

resource "vault_kv_secret_v2" "account_static_key" {
  mount               = var.secret_mount
  name                = "${var.secret_name}/account_static_key"
  data_json = jsonencode(
    {
        id = yandex_iam_service_account_static_access_key.static_key.id
        cloud_id = yandex_resourcemanager_folder.folder.cloud_id
        folder_id = yandex_resourcemanager_folder.folder.folder_id
        service_account_id = yandex_iam_service_account.account.service_account_id
        access_key = yandex_iam_service_account_static_access_key.static_key.access_key
        secret_key = yandex_iam_service_account_static_access_key.static_key.secret_key
    }
  )
}

resource "vault_kv_secret_v2" "account_iam_key" {
  mount               = var.secret_mount
  name                = "${var.secret_name}/account_iam_key"
  data_json = jsonencode(
    {
        id = yandex_iam_service_account_key.iam_key.id
        cloud_id = yandex_resourcemanager_folder.folder.cloud_id
        folder_id = yandex_resourcemanager_folder.folder.folder_id
        public_key = yandex_iam_service_account_key.iam_key.public_key
        private_key = yandex_iam_service_account_key.iam_key.private_key
        service_account_key = jsonencode({
            id = yandex_iam_service_account_key.iam_key.id
            service_account_id = yandex_iam_service_account.account.service_account_id
            created_at = yandex_iam_service_account_key.iam_key.created_at
            key_algorithm = yandex_iam_service_account_key.iam_key.key_algorithm
            public_key = yandex_iam_service_account_key.iam_key.public_key
            private_key = yandex_iam_service_account_key.iam_key.private_key
        })
    }
  )
}
