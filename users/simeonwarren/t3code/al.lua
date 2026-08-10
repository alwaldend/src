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
        vault_secret = "alwaldend.com/vault1/approles/user_simeonwarren/tf_backend/t3code_tf_setup",
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
    name = "wireguard",
    plugin = "injector",
    labels = { tf = "setup", ansible = "1" },
    data = {
        res = {
            {
                name = "wireguard",
                kv = {
                    path = "alwaldend.com/vault1/approles/user_simeonwarren/t3code/wireguard",
                    mount = "secrets"
                }
            },
            {
                name = "TF_VAR_wg_public_keys",
                deps = { "wireguard" },
                env = {
                    value = "{{ .Last.Data.wg_public_keys | to_json }}",
                }
            },
            {
                name = "TF_VAR_wg_private_keys",
                deps = { "wireguard" },
                env = {
                    value = "{{ .Last.Data.wg_private_keys | to_json }}",
                }
            },
            {
                name = "TF_VAR_wg_preshared_keys",
                deps = { "wireguard" },
                env = {
                    value = "{{ .Last.Data.wg_preshared_keys | to_json }}",
                }
            },
        },
    },
})

infra.mikrotik({
    path = "alwaldend.com/vault1/approles/user_simeonwarren/mikrotik",
    host = "https://router1.dc1.alwaldend.com",
    labels = { tf = "setup" },
})

lib.plugin_call({
    name = "http_proxies",
    plugin = "injector",
    labels = { ansible = "1" },
    data = {
        res = {
            {
                name = "http_proxies",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_threexui/http_proxies",
                    mount = "secrets"
                }
            },
            {
                name = "HTTP_PROXIES",
                deps = { "http_proxies" },
                env = {
                    value = "{{ .Last.Data.http_proxies | to_json }}",
                }
            },
        },
    },
})
