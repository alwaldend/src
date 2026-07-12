local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_third_party",
    },
})

lib.plugin_call({
    name = "forgejo_login",
    plugin = "forgejo_login",
    labels = { forgejo = "1" },
})
