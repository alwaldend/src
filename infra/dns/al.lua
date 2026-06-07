local lib = require("al_lib")
lib.vault_auth({
    name = "default",
    approle = { name = "src_infra_dns" }
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
                    mount = "secrets"
                }
            },
            {
                name = "mikrotik",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_dns/mikrotik",
                    mount = "secrets"
                }
            },
            {
                name = "creds",
                deps = { "cloudflare", "mikrotik" },
                file = {
                    from_file = "infra/dns/creds.json.tpl",
                }
            },
            {
                name = "DNSCONTROL_CREDS",
                deps = {"creds"},
                env = {
                    value = "{{ index .Last.Files 0 }}",
                }
            },
        }
    }
})
