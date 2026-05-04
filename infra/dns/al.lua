local lib = require("al_lib")
lib.auth({
    name = "default",
    approle = { name = "src_infra_dns" }
})
lib.secret({
    name = "dns_token",
    kv = { path = "cloudflare.com/dns_token", mount = "secrets" }
})
lib.env({
    name = "CLOUDFLARE_ACCOUNT_ID",
    secrets = {"dns_token"},
    labels = { dns = "1" },
    value = "{{ .Secret.cloudflare_account_id }}"
})
lib.env({
    name = "CLOUDFLARE_API_TOKEN",
    secrets = {"dns_token"},
    labels = { dns = "1" } ,
    value = "{{ .Secret.cloudflare_api_token }}"
})
