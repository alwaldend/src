local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.auth({
    name = "default",
    approle = {
        name = "src_infra_dc1_vault",
    },
})

lib.secret({
    name = "vault_certs",
    kv = {
        path = "alwaldend.com/vault1/approles/src_infra_dc1_vault/vault_certs",
        mount = "secrets",
    },
    labels = { ansible = "1" }
})

infra.tf_backend({
    path = "alwaldend.com/vault1/approles/src_infra_dc1_vault/bucket",
    labels = { tf = "1" },
})

infra.yc_auth({
    path = "yandex.cloud/org1/folders/src-infra-dc1-vault/account_iam_key",
    labels = { tf = "1" },
})

infra.yc_bucket_auth({
    path = "yandex.cloud/org1/folders/src-infra-dc1-vault/account_static_key",
    labels = { tf = "1", rclone = "1" },
})

infra.yc_account({
    path = "yandex.cloud/org1/folders/src-infra-dc1-vault/account",
    labels = { tf = "1" },
})

infra.rclone_config({
    path = "alwaldend.com/vault1/approles/src_infra_dc1_vault/backup_bucket",
    labels = { rclone = "1" },
})
