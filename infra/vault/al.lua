local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_dc1_vault",
    },
})

infra.ansible_keys({
    name = "ansible_keys",
    labels = { ansible = "1" },
    vault_ssh = {
        backend = "ssh/clients/sign/admins",
        ttl = 60 * 60 * 2
    }
})

lib.plugin_call({
    name = "tf_backend",
    plugin = "tf_backend",
    labels = { tf = "setup" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_dc1_vault/tf_backend/tf_setup",
        vault_secret_mount = "secrets"
    },
})

lib.plugin_call({
    name = "tf_backend",
    plugin = "tf_backend",
    labels = { tf = "main" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_dc1_vault/tf_backend/main",
        vault_secret_mount = "secrets"
    },
})

infra.yc_auth({
    path = "yandex.cloud/org1/folders/src-infra-dc1-vault/account_iam_key",
    labels = { tf = "main" },
})

infra.yc_bucket_auth({
    path = "yandex.cloud/org1/folders/src-infra-dc1-vault/account_static_key",
    labels = { tf = "main" },
})

infra.yc_account({
    path = "yandex.cloud/org1/folders/src-infra-dc1-vault/account",
    labels = { tf = "main" },
})

infra.rclone_config({
    path = "alwaldend.com/vault1/approles/src_infra_dc1_vault/backup_bucket",
    labels = { backup = "1" },
})

lib.plugin_call({
    name = "unseal",
    plugin = "injector",
    labels = { unseal = "1" },
    data = {
        res = {
            {
                name = "unseal",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_dc1_vault/unseal",
                    mount = "secrets"
                }
            },
            {
                name = "UNSEAL_KEY",
                deps = { "unseal" },
                env = {
                    value = "{{ .Last.Data.encrypted_key }}",
                }
            },
        }
    }
})
