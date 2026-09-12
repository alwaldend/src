local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_projects_alwaldend_com",
    },
})

lib.plugin_call({
    name = "tf_backend_tf",
    plugin = "tf_backend",
    labels = { tf = "main", dns = "1" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_projects_alwaldend_com/tf_backend/tf",
        vault_secret_mount = "secrets",
    },
})

lib.plugin_call({
    name = "pve_login",
    plugin = "pve_login",
    labels = { tf = "main" },
})

infra.pve_provider_inputs({
    labels = { dns = "1" },
})

infra.dns({
    labels = { tf = "main", dns = "1" },
    dc1 = false,
})
