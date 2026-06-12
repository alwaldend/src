local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_dc1_forgejo1",
    },
})

lib.plugin_call({
    name = "tf_backend_tf_setup",
    plugin = "tf_backend",
    labels = { tf_setup = "1" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_dc1_forgejo1/tf_backend/tf_setup",
        vault_secret_mount = "secrets"
    },
})

lib.plugin_call({
    name = "tf_backend_tf",
    plugin = "tf_backend",
    labels = { tf_main = "1" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_dc1_forgejo1/tf_backend/tf",
        vault_secret_mount = "secrets"
    },
})

lib.plugin_call({
    name = "pve_login",
    plugin = "pve_login",
    labels = { tf_setup = "1" },
})
