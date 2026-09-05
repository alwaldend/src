local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_xcp_ng",
    },
})

lib.plugin_call({
    name = "tf_backend",
    plugin = "tf_backend",
    labels = { tf = "main" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_xcp_ng/tf_backend",
        vault_secret_mount = "secrets",
    },
})
