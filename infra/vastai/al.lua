local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_vastai",
    },
})

lib.plugin_call({
    name = "tf_backend_tf",
    plugin = "tf_backend",
    labels = { tf = "main" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_vastai/tf_backend/tf",
        vault_secret_mount = "secrets",
    },
})

lib.plugin_call({
    name = "vastai",
    plugin = "injector",
    labels = { vastai = "1" },
    data = {
        res = {
            {
                name = "vastai",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_vastai/vastai",
                    mount = "secrets",
                },
            },
            {
                name = "VAST_API_KEY",
                deps = { "vastai" },
                env = {
                    value = "{{ .Last.Data.api_key }}",
                },
            },
        },
    },
})
