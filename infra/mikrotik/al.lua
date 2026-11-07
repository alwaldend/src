local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = { name = "src_infra_mikrotik" },
})

lib.plugin_call({
    name = "tf_backend",
    plugin = "tf_backend",
    labels = { tf = "main" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_mikrotik/tf_backend/tf",
        vault_secret_mount = "secrets",
    },
})

infra.dns({
    labels = { tf = "main" },
    dc1 = true,
    global = true,
})

lib.plugin_call({
    name = "firewall_backend",
    plugin = "tf_backend",
    labels = { tf = "firewall" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_mikrotik/tf_backend/firewall",
        vault_secret_mount = "secrets",
    },
})

lib.plugin_call({
    name = "firewall_routeros",
    plugin = "injector",
    labels = { tf = "firewall" },
    data = {
        res = {
            {
                name = "firewall_routeros",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_ingress/mikrotik",
                    mount = "secrets",
                },
            },
            {
                name = "TF_VAR_routeros_hosturl",
                deps = { "firewall_routeros" },
                env = { value = infra.routeros_hosturl },
            },
            {
                name = "TF_VAR_routeros_username",
                deps = { "firewall_routeros" },
                env = { value = "{{ .Last.Data.mikrotik_username }}" },
            },
            {
                name = "TF_VAR_routeros_password",
                deps = { "firewall_routeros" },
                env = { value = "{{ .Last.Data.mikrotik_password }}" },
            },
        },
    },
})
