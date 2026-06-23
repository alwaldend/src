local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    token_helper = nil
})

lib.vault_auth({
    name = "token_helper",
    token_helper = nil
})

lib.vault_auth({
    name = "no_auth",
    no_auth = true
})

lib.vault_conn({
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

lib.plugin_call({
    name = "env_default_no_auth",
    plugin = "injector",
    labels = { vault_env = "default_no_auth" },
    data = {
        res = {
            {
                name = "env_default_no_auth",
                vault_env = {
                    conn = "default",
                    auth = "no_auth",
                },
            },
        }
    },
})

lib.plugin_call({
    name = "env_default_default",
    plugin = "injector",
    labels = { vault_env = "default_default" },
    data = {
        res = {
            {
                name = "env_default_default",
                vault_env = {
                    conn = "default",
                    auth = "default",
                },
            },
        }
    },
})

lib.plugin({
    name = "injector",
    bin = "com_alwaldend_src/tools/vault/injector/injector_/injector",
})

lib.plugin({
    name = "forgejo_login",
    bin = "com_alwaldend_src/tools/vault/forgejo_login/forgejo_login_/forgejo_login",
    labels = {
        forgejo_login = "1"
    },
    data = {
        forgejo_url = "https://forgejo1.dc1.alwaldend.com",
    },
})

lib.plugin({
    name = "pve_login",
    bin = "com_alwaldend_src/tools/vault/pve_login/pve_login_/pve_login",
    labels = {
        pve_login = "1"
    },
    data = {
        pve_base_url = "https://host1.pve1.dc1.alwaldend.com:8006",
        pve_redirect_url = "https://host1.pve1.dc1.alwaldend.com:8006",
        pve_realm = "src_infra_dc1_vault",
    },
})

lib.plugin({
    name = "tf_backend",
    bin = "com_alwaldend_src/tools/vault/tf_backend/tf_backend_/tf_backend",
})
