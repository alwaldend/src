local lib = require("al_lib")

for _, stage in ipairs({ "rentals", "rental_test" }) do
    lib.plugin_call({
        name = "tf_backend_" .. stage,
        plugin = "tf_backend",
        labels = { tf = stage },
        data = {
            vault_secret = "alwaldend.com/vault1/approles/src_infra_vastai/tf_backend/"
                .. stage,
            vault_secret_mount = "secrets",
        },
    })
end

lib.plugin_call({
    name = "vastai_rentals",
    plugin = "injector",
    labels = { vastai = "rentals" },
    data = {
        res = {
            {
                name = "vastai_rentals",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_vastai/vastai/alwaldend",
                    mount = "secrets",
                },
            },
            {
                name = "VAST_API_KEY",
                deps = { "vastai_rentals" },
                env = { value = "{{ .Last.Data.api_key }}" },
            },
            {
                name = "TF_VAR_alwaldend_api_key",
                deps = { "vastai_rentals" },
                env = { value = "{{ .Last.Data.api_key }}" },
            },
        },
    },
})

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
                name = "env_component",
                vault_env = { conn = "default", auth = "default" },
            },
            {
                name = "TF_VAR_component_vault_token",
                deps = { "env_component" },
                env = { value = "{{ .Last.Env.VAULT_TOKEN }}" },
            },
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
