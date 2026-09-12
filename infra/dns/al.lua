local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = { name = "src_infra_dns" },
})

lib.plugin_call({
    name = "tf_backend_tf",
    plugin = "tf_backend",
    labels = { tf = "1" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_dns/tf_backend/tf",
        vault_secret_mount = "secrets",
    },
})

infra.dns({
    labels = { tf = "1" },
    dc1 = true,
})
