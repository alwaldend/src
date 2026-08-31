local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = { name = "src_infra_dns" },
})

lib.plugin_call({
    name = "tf_backend_tf",
    plugin = "tf_backend",
    labels = { tf = "1" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_dns/tf_backend/tf",
        vault_secret_mount = "secrets",
    },
})

lib.plugin_call({
    name = "cloudflare",
    plugin = "injector",
    labels = { tf = "1" },
    data = {
        res = {
            {
                name = "cloudflare",
                kv = {
                    path = "cloudflare.com/dns_token",
                    mount = "secrets",
                },
            },
            {
                name = "CLOUDFLARE_ACCOUNT_ID",
                deps = { "cloudflare" },
                env = {
                    value = "{{ .Last.Data.cloudflare_account_id }}",
                },
            },
            {
                name = "CLOUDFLARE_API_TOKEN",
                deps = { "cloudflare" },
                env = {
                    value = "{{ .Last.Data.cloudflare_api_token }}",
                },
            },
        },
    },
})

lib.plugin_call({
    name = "dns",
    plugin = "injector",
    labels = { dns = "1" },
    data = {
        res = {
            {
                name = "cloudflare",
                kv = {
                    path = "cloudflare.com/dns_token",
                    mount = "secrets",
                },
            },
            {
                name = "mikrotik",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_dns/mikrotik",
                    mount = "secrets",
                },
            },
            {
                name = "creds",
                deps = { "cloudflare", "mikrotik" },
                file = {
                    from_file = "infra/dns/creds.json.tpl",
                },
            },
            {
                name = "DNSCONTROL_CREDS",
                deps = { "creds" },
                env = {
                    value = "{{ index .Last.Files 0 }}",
                },
            },
        },
    },
})
