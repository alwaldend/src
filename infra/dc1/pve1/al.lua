local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_dc1_pve1",
    },
})

lib.plugin_call({
    name = "tf_backend",
    plugin = "tf_backend",
    labels = { tf = "1" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_dc1_pve1/tf_backend",
        vault_secret_mount = "secrets"
    },
})

infra.yc_auth({
    path = "yandex.cloud/org1/folders/src-infra-dc1-pve1/account_iam_key",
    labels = { tf = "1" },
})

infra.yc_bucket_auth({
    path = "yandex.cloud/org1/folders/src-infra-dc1-pve1/account_static_key",
    labels = { tf = "1" },
})

infra.yc_account({
    path = "yandex.cloud/org1/folders/src-infra-dc1-pve1/account",
    labels = { tf = "1" },
})
