local lib = require("al_lib")
local infra = require("infra.al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_harbor",
    },
})

infra.ansible_keys({
    name = "ansible_keys",
    labels = { ansible = "1" },
    vault_ssh = {
        backend = "ssh/clients/sign/admins",
        ttl = 60 * 60 * 2
    }
})

lib.plugin_call({
    name = "tf_backend_tf_setup",
    plugin = "tf_backend",
    labels = { tf = "setup" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_harbor/tf_backend/tf_setup",
        vault_secret_mount = "secrets"
    },
})

lib.plugin_call({
    name = "tf_backend_tf",
    plugin = "tf_backend",
    labels = { tf = "main" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_harbor/tf_backend/tf",
        vault_secret_mount = "secrets"
    },
})

lib.plugin_call({
    name = "pve_login",
    plugin = "pve_login",
    labels = { tf = "setup" },
})

lib.plugin_call({
    name = "config",
    plugin = "injector",
    labels = { ansible = "1" },
    data = {
        res = {
            {
                name = "config",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_harbor/config",
                    mount = "secrets"
                }
            },
            {
                name = "HARBOR_ADMIN_PASSWORD",
                deps = { "config" },
                env = {
                    value = "{{ .Last.Data.admin_password }}",
                }
            },
            {
                name = "HARBOR_DB_PASSWORD",
                deps = { "config" },
                env = {
                    value = "{{ .Last.Data.db_password }}",
                }
            },
        }
    }
})
