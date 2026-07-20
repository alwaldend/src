local infra = require("infra.al_lib")
local lib = require("al_lib")

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
    name = "tf_backend",
    plugin = "tf_backend",
    labels = { tf = "1" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/user_simeonwarren/tf_backend/main",
        vault_secret_mount = "secrets"
    },
})

infra.yc_auth({
    path = "yandex.cloud/org1/folders/user-simeonwarren/account_iam_key",
    labels = { tf = "1" },
})

lib.plugin_call({
    name = "oidc_client",
    plugin = "injector",
    labels = { ansible = "1" },
    data = {
        res =  {
            {
                name = "oidc_client",
                op = {
                    method = "read",
                    path = "identity/oidc/client/users_simeonwarren_provider"
                }
            },
            {
                name = "OIDC_CLIENT_ID",
                deps = { "oidc_client" },
                env = {
                    value = "{{ .Last.Data.client_id }}",
                }
            },
            {
                name = "OIDC_CLIENT_SECRET",
                deps = { "oidc_client" },
                env = {
                    value = "{{ .Last.Data.client_secret }}",
                }
            },
        }
    }
})

lib.plugin_call({
    name = "pve_login",
    plugin = "pve_login",
    labels = { tf = "1" },
})

infra.yc_account({
    path = "yandex.cloud/org1/folders/user-simeonwarren/account",
    labels = { tf = "1" },
})

infra.yc_bucket_auth({
    path = "yandex.cloud/org1/folders/user-simeonwarren/account_static_key",
    labels = { tf = "1" },
})

lib.plugin_call({
    name = "opencode",
    plugin = "injector",
    labels = { opencode = "1" },
    data = {
        res = {
            {
                name = "opencode",
                kv = {
                    path = "alwaldend.com/vault1/approles/user_simeonwarren/opencode",
                    mount = "secrets"
                }
            },
            {
                name = "OPENCODE_YANDEX_CLOUD_TOKEN",
                deps = { "opencode" },
                env = {
                    value = "{{ .Last.Data.secret_key }}",
                }
            },
            {
                name = "OPENCODE_YANDEX_CLOUD_FOLDER",
                deps = { "opencode" },
                env = {
                    value = "{{ .Last.Data.folder_id }}",
                }
            },
        }
    }
})
