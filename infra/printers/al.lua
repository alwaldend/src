local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = { name = "src_infra_printers" },
})

lib.plugin_call({
    name = "tf_backend",
    plugin = "tf_backend",
    labels = { tf = "dns" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_printers/tf_backend/dns",
        vault_secret_mount = "secrets",
    },
})

infra.dns({
    labels = { tf = "dns" },
    dc1 = true,
    global = true,
})

lib.plugin_call({
    name = "firewall_backend",
    plugin = "tf_backend",
    labels = { tf = "main" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_printers/tf_backend/tf",
        vault_secret_mount = "secrets",
    },
})

infra.mikrotik({
    name = "firewall_routeros",
    labels = { tf = "main" },
    path = "alwaldend.com/vault1/approles/src_infra_mikrotik/printers",
    host = infra.routeros_hosturl,
})
