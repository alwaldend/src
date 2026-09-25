local lib = require("al_lib")

lib.plugin_call({
    name = "download_approle_backend",
    plugin = "tf_backend",
    labels = { tf = "approle_download" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_dc1_vault/tf_backend/approles/src_infra_download",
        vault_secret_mount = "secrets",
    },
})
