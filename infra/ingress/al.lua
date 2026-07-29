local lib = require("al_lib")
local infra = require("infra.al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_ingress",
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
    name = "tf_backend_tf",
    plugin = "tf_backend",
    labels = { tf = "main" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_ingress/tf_backend/tf",
        vault_secret_mount = "secrets"
    },
})

infra.yc_auth({
    path = "yandex.cloud/org1/folders/src-infra-ingress/account_iam_key",
    labels = { tf = "main" },
})

infra.yc_account({
    path = "yandex.cloud/org1/folders/src-infra-ingress/account",
    labels = { tf = "main" },
})

lib.plugin_call({
    name = "image",
    plugin = "injector",
    labels = { tf = "main" },
    data = {
        res =  {
            {
                name = "image",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_ingress/image",
                    mount = "secrets",
                }
            },
            {
                name = "TF_VAR_image_signed_url",
                deps = { "image" },
                env = {
                    value = "{{ .Last.Data.image_signed_url }}",
                }
            },
        }
    }
})

lib.plugin_call({
    name = "mikrotik",
    plugin = "injector",
    labels = { tf = "main" },
    data = {
        res = {
            {
                name = "mikrotik",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_ingress/mikrotik",
                    mount = "secrets"
                }
            },
            {
                name = "MIKROTIK_HOST",
                deps = { "mikrotik" },
                env = {
                    value = "https://router1.dc1.alwaldend.com",
                }
            },
            {
                name = "MIKROTIK_USER",
                deps = { "mikrotik" },
                env = {
                    value = "{{ .Last.Data.mikrotik_username }}",
                }
            },
            {
                name = "MIKROTIK_PASSWORD",
                deps = { "mikrotik" },
                env = {
                    value = "{{ .Last.Data.mikrotik_password }}",
                }
            },
            {
                name = "wireguard",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_ingress/wireguard",
                    mount = "secrets"
                }
            },
            {
                name = "TF_VAR_wg_public_key",
                deps = { "wireguard" },
                env = {
                    value = "{{ .Last.Data.wg_node_public_key }}",
                }
            },
            {
                name = "TF_VAR_wg_preshared_key",
                deps = { "wireguard" },
                env = {
                    value = "{{ .Last.Data.wg_preshared_key }}",
                }
            },
            {
                name = "TF_VAR_wg_interface_private_key",
                deps = { "wireguard" },
                env = {
                    value = "{{ .Last.Data.wg_router_private_key }}",
                }
            },
        }
    }
})
