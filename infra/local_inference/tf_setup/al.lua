local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = { name = "src_infra_local_inference" },
})

lib.plugin_call({
    name = "tf_backend",
    plugin = "tf_backend",
    labels = { tf = "setup" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_local_inference/tf_backend/tf_setup",
        vault_secret_mount = "secrets",
    },
})

infra.yc_auth({
    path = "yandex.cloud/org1/folders/src-infra-local-inference/account_iam_key",
    labels = { tf = "setup" },
})
