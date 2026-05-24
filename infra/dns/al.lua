local lib = require("al_lib")
lib.auth({
    name = "default",
    approle = { name = "src_infra_dns" }
})
lib.secret({
    name = "cloudflare",
    kv = { path = "cloudflare.com/dns_token", mount = "secrets" }
})
lib.secret({
    name = "mikrotik",
    kv = { path = "alwaldend.com/vault1/approles/src_infra_dns/mikrotik", mount = "secrets" }
})
lib.file({
    name = "creds",
    secrets = { "cloudflare", "mikrotik" },
    from_file = "infra/dns/creds.json.tpl",
})
lib.env({
    name = "DNSCONTROL_CREDS",
    files = {"creds"},
    labels = { dns = "1" },
    value = "{{ .File.Path }}"
})
