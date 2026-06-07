local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = { name = "src_infra_yandex_cloud_org1" },
})

lib.plugin_call({
    name = "tf_backend",
    plugin = "tf_backend",
    labels = { tf = "1" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_yandex_cloud_org1/tf_backend/main",
        vault_secret_mount = "secrets"
    },
})

infra.yc_auth({
    path = "yandex.cloud/org1/folders/src-infra-yandex-cloud-org1/account_iam_key",
    labels = { tf = "1" },
})

infra.yc_account({
    path = "yandex.cloud/org1/folders/src-infra-yandex-cloud-org1/account",
    labels = { tf = "1" },
})

infra.yc_bucket_auth({
    path = "yandex.cloud/org1/folders/src-infra-yandex-cloud-org1/account_static_key",
    labels = { tf = "1" },
})
