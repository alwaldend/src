local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_projects_x_article_uploader",
    },
})

-- The publication credential is an OAuth 1.0a user-context credential held at
-- the uploader's own path under its AppRole subtree, like the other components'
-- credentials. Only the reference is checked in; the values are written by an
-- authorized operator and injected at run time. The shared AppRole module's
-- own-subtree policy grants the read.
lib.plugin_call({
    name = "x_api",
    plugin = "injector",
    labels = { x_api = "1" },
    data = {
        res = {
            {
                name = "x_api",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_projects_x_article_uploader/oauth1",
                    mount = "secrets",
                },
            },
            {
                name = "X_API_KEY",
                deps = { "x_api" },
                env = {
                    value = "{{ .Last.Data.api_key }}",
                },
            },
            {
                name = "X_API_SECRET",
                deps = { "x_api" },
                env = {
                    value = "{{ .Last.Data.api_secret }}",
                },
            },
            {
                name = "X_ACCESS_TOKEN",
                deps = { "x_api" },
                env = {
                    value = "{{ .Last.Data.access_token }}",
                },
            },
            {
                name = "X_ACCESS_TOKEN_SECRET",
                deps = { "x_api" },
                env = {
                    value = "{{ .Last.Data.access_token_secret }}",
                },
            },
        },
    },
})
