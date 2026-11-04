local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = { name = "src_projects_vastai_terraform_provider" },
})

lib.plugin_call({
    name = "tf_backend_test",
    plugin = "tf_backend",
    labels = { tf = "test" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_projects_vastai_terraform_provider/tf_backend/test",
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
                    path = "alwaldend.com/vault1/approles/src_projects_vastai_terraform_provider/vastai",
                    mount = "secrets",
                },
            },
            {
                name = "VAST_API_KEY",
                deps = { "vastai" },
                env = { value = "{{ .Last.Data.api_key }}" },
            },
        },
    },
})
