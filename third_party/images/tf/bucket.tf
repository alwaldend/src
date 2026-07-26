resource "yandex_kms_symmetric_key" "images" {
  name              = "images"
  folder_id         = var.folder_id
  labels            = local.labels
  description       = "Key for the bucket"
  default_algorithm = "AES_128"
  rotation_period   = "8760h" // equal to 1 year
}

resource "yandex_storage_bucket" "images" {
  bucket                = "com-alwaldend-src-third-party-images"
  folder_id             = var.folder_id
  tags                  = local.labels
  max_size              = (1024 * 1024 * 1024) * 1000
  default_storage_class = "COLD"
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = yandex_kms_symmetric_key.images.id
        sse_algorithm     = "aws:kms"
      }
    }
  }
}

resource "yandex_storage_bucket_iam_binding" "images_writer" {
  bucket = yandex_storage_bucket.images.bucket
  role   = "storage.uploader"
  members = [
    "serviceAccount:${var.service_account_id}",
  ]
}

resource "yandex_storage_bucket_iam_binding" "images_reader" {
  bucket = yandex_storage_bucket.images.bucket
  role   = "storage.viewer"
  members = [
    "system:group:organization:bpftkdqqdr0sqec1ma0j:users",
  ]
}
