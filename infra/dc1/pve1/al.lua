local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.auth({
    name = "default",
    approle = {
        name = "src_infra_dc1_pve1",
    },
})

lib.plugin({
    name = "vault_tf_backend_local",
    extends = { "vault_tf_backend" },
    labels = { tf = "1" },
    config = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_dc1_pve1/tf_backend"
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
