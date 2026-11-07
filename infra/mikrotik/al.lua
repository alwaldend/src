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
    name = "accounts_backend",
    plugin = "tf_backend",
    labels = { tf = "accounts" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_mikrotik/tf_backend/accounts",
        vault_secret_mount = "secrets",
    },
})

infra.mikrotik({
    name = "accounts_routeros",
    labels = { tf = "accounts" },
    path = "alwaldend.com/vault1/approles/src_infra_mikrotik/routeros_admin",
    host = infra.routeros_hosturl,
})

lib.plugin_call({
    name = "accounts_vault_env",
    plugin = "injector",
    labels = { tf = "accounts" },
    data = {
        res = {
            {
                name = "accounts_vault_env",
                vault_env = { conn = "default", auth = "default" },
            },
        },
    },
})
