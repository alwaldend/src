local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_vastai",
    },
})

lib.vault_auth({
    name = "provider_project",
    approle = { name = "src_projects_vastai_terraform_provider" },
})

lib.plugin_call({
    name = "env_provider_project",
    plugin = "injector",
    labels = { vault_env = "provider_project" },
    data = {
        res = {
            {
                name = "env_provider_project",
                vault_env = { conn = "default", auth = "provider_project" },
            },
        },
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
