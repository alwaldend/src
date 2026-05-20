local lib = require("al_lib")

lib.vault({
    name = "default",
    config = {
        address = "https://vault.dc1.alwaldend.com:8200"
    },
    tls = {
        ca_cert = "data/ssl/alwaldend.com/root_ca.crt",
        -- client_cert = "${HOME}/.al/client_cert/host.crt",
        -- client_key = "${HOME}/.al/client_cert/host.key",
    }
})

lib.env({
    name = "VAULT_SKIP_VERIFY",
    value = "1",
})

lib.auth({
    name = "default",
    token_helper = nil
})

lib.auth({
    name = "token_helper",
    token_helper = nil
})

lib.auth({
    name = "no_auth",
    no_auth = true
})
