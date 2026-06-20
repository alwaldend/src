resource "yandex_iam_service_account" "opencode" {
  name        = "opencode"
  folder_id   = var.folder_id
  description = "Opencode account"
}

resource "yandex_resourcemanager_folder_iam_member" "opencode" {
  role      = "ai.languageModels.user"
  folder_id = var.folder_id
  member    = "serviceAccount:${yandex_iam_service_account.opencode.id}"
}


resource "yandex_iam_service_account_api_key" "opencode" {
  lifecycle {
    ignore_changes = [scope]
  }
  service_account_id = yandex_iam_service_account.opencode.id
  description        = "Api key for Opencode"
  scopes             = ["yc.ai.languageModels.execute"]
}

resource "vault_kv_secret_v2" "opencode" {
  mount = local.kv_mount
  name  = "${local.kv_path}/opencode"
  data_json = jsonencode(
    {
      folder_id  = var.folder_id
      secret_key = yandex_iam_service_account_api_key.opencode.secret_key
    }
  )
}
