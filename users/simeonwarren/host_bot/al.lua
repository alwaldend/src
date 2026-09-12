local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = { name = "user_simeonwarren" },
})

infra.ansible_keys({
    name = "ansible_keys",
    labels = { ansible = "1" },
    vault_ssh = {
        backend = "ssh/clients/sign/admins",
        ttl = 60 * 60 * 2,
    },
})

lib.plugin_call({
    name = "openhands",
    plugin = "injector",
    labels = { ansible = "1" },
    data = {
        res = {
            {
                name = "openhands",
                kv = {
                    path = "alwaldend.com/vault1/approles/user_simeonwarren/openhands",
                    mount = "secrets",
                },
            },
            {
                name = "OPENHANDS_SERVER_SECRET_KEY",
                deps = { "openhands" },
                env = {
                    value = "{{ .Last.Data.secret_key }}",
                },
            },
        },
    },
})
