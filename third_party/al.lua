local lib = require("al_lib")
local infra = require("infra.al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_third_party",
    },
})

lib.plugin_call({
    name = "forgejo_login",
    plugin = "forgejo_login",
    labels = { forgejo = "1" },
})

lib.plugin_call({
    name = "tf_backend_tf",
    plugin = "tf_backend",
    labels = { tf = "main" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_third_party/tf_backend/tf",
        vault_secret_mount = "secrets"
    },
})

infra.yc_bucket_auth({
    path = "yandex.cloud/org1/folders/src-third-party/account_static_key",
    labels = { tf = "main", rclone = "1" },
})


infra.yc_auth({
    path = "yandex.cloud/org1/folders/src-third-party/account_iam_key",
    labels = { tf = "main" },
})

infra.yc_account({
    path = "yandex.cloud/org1/folders/src-third-party/account",
    labels = { tf = "main" },
})
