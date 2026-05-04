local lib = require("al_lib")

lib.vault({
    name = "default",
    config = {
        address = "https://vault.dc1.alwaldend.com:8200"
    },
    tls = {
        ca_cert = "data/ssl/alwaldend.com_root_ca.crt"
    }
})

lib.auth({
    name = "default",
    token_helper = nil
})
