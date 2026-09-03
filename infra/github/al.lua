local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_github",
    },
})

lib.plugin_call({
    name = "tf_backend_tf",
    plugin = "tf_backend",
    labels = { tf = "main" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_github/tf_backend/tf",
        vault_secret_mount = "secrets",
    },
})

lib.plugin_call({
    name = "github",
    plugin = "injector",
    labels = { tf = "main" },
    data = {
        res = {
            {
                name = "github",
                kv = {
                    path = "github.com/pages_token",
                    mount = "secrets",
                },
            },
            {
                name = "GITHUB_TOKEN",
                deps = { "github" },
                env = {
                    value = "{{ .Last.Data.token }}",
                },
            },
        },
    },
})
