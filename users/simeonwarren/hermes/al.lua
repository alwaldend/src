local lib = require("al_lib")
local infra = require("infra.al_lib")

lib.vault_auth({
    name = "default",
    approle = { name = "user_simeonwarren" },
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
    name = "tf_backend_tf_setup",
    plugin = "tf_backend",
    labels = { tf = "setup" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/user_simeonwarren/tf_backend/hermes_tf_setup",
        vault_secret_mount = "secrets"
    },
})

infra.yc_auth({
    path = "yandex.cloud/org1/folders/user-simeonwarren/account_iam_key",
    labels = { tf = "setup" },
})

infra.yc_account({
    path = "yandex.cloud/org1/folders/user-simeonwarren/account",
    labels = { tf = "setup" },
})

lib.plugin_call({
    name = "hermes_config",
    plugin = "injector",
    labels = { ansible = "1" },
    data = {
        res = {
            {
                name = "config",
                kv = {
                    path = "alwaldend.com/vault1/approles/user_simeonwarren/hermes_config",
                    mount = "secrets"
                }
            },
            {
                name = "HERMES_ENV",
                deps = { "config" },
                env = {
                    value = '{{ to_json_indent .Last.Data.env "" "  " }}',
                }
            },
        }
    }
})
