local lib = require("al_lib")
lib.auth({
    name = "default",
    approle = { name = "src_infra_dns" }
})
lib.secret({
    name = "dns_token",
    kv = { path = "cloudflare.com/dns_token", mount = "secrets" }
})
lib.file({
    name = "creds",
    secrets = {"dns_token"},
    from_file = "infra/dns/creds.json.tpl",
})
lib.env({
    name = "DNSCONTROL_CREDS",
    files = {"creds"},
    labels = { dns = "1" },
    value = "{{ .File.Path }}"
})
