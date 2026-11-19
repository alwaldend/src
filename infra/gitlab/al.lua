local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_gitlab",
    },
})

lib.plugin_call({
    name = "tf_backend_tf",
    plugin = "tf_backend",
    labels = { tf = "main" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_gitlab/tf_backend/tf",
        vault_secret_mount = "secrets",
    },
})

lib.plugin_call({
    name = "gitlab",
    plugin = "injector",
    labels = { tf = "main" },
    data = {
        res = {
            {
                name = "gitlab",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_gitlab/gitlab",
                    mount = "secrets",
                },
            },
            {
                name = "GITLAB_TOKEN",
                deps = { "gitlab" },
                env = {
                    value = "{{ .Last.Data.gitlab_token }}",
                },
            },
        },
    },
})
