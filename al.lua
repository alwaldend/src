local lib = require("al_lib")

lib.vault({
    name = "default",
    config = {
        address = "https://vault.dc1.alwaldend.com:8200"
    },
    tls = {
        ca_cert = "infra/dc1/vault/tf/output/pki_ca_servers.crt",
        -- client_cert = "${HOME}/.al/client_cert/host.crt",
        -- client_key = "${HOME}/.al/client_cert/host.key",
    }
})

lib.plugin({
    name = "pve_login",
    bin = "com_alwaldend_src/tools/vault/pve_login/pve_login_/pve_login",
    config = {
        pve_base_url = { value_string = "https://host1.pve1.dc1.alwaldend.com:8006" },
        pve_redirect_url = { value_string = "https://host1.pve1.dc1.alwaldend.com:8006" } ,
        pve_realm = { value_string = "src_infra_dc1_vault" },
    },
})

lib.plugin({
    name = "vault_tf_backend",
    bin = "com_alwaldend_src/tools/vault/tf_backend/tf_backend_/tf_backend",
    config = {
        vault_secret_mount = { value_string = "secrets" }
    }
})

lib.env({
    name = "VAULT_SKIP_VERIFY",
    value = "1",
    labels = { insecure = "1" }
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
