local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.plugin({
    name = "xo_login",
    bin = "com_alwaldend_src/infra/xcp_ng/cmd/xo_login/xo_login_/xo_login",
    labels = { xoa_login = "1" },
    data = {
        xoa_url = "https://xoa.xcp-ng.alwaldend.com",
        discovery_url = "https://vault.alwaldend.com:8200/v1/identity/oidc/provider/src_infra_xcp_ng_provider/.well-known/openid-configuration",
    },
})

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_forgejo_runner",
    },
})

-- Registration runs on the controller with an invocation-scoped Forgejo
-- administrator token. This identity is never installed on the runner VM.
lib.vault_auth({
    name = "forgejo_registration",
    approle = { name = "src_infra_dc1_forgejo1" },
})

lib.plugin_call({
    name = "forgejo_registration",
    plugin = "forgejo_login",
    labels = { registration = "1" },
    data = { vault_auth = "forgejo_registration" },
})

infra.ansible_keys({
    name = "ansible_keys",
    labels = { ansible = "1" },
    vault_ssh = {
        backend = "ssh/clients/sign/admins",
        ttl = 60 * 60 * 2,
    },
})

lib.plugin_call({
    name = "tf_backend_tf_setup",
    plugin = "tf_backend",
    labels = { tf = "setup", dns = "1", runner = "1" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_forgejo_runner/tf_backend/tf_setup",
        vault_secret_mount = "secrets",
    },
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
                    path = "alwaldend.com/vault1/approles/src_infra_forgejo_runner/config",
                    mount = "secrets",
                },
            },
            {
                name = "FORGEJO_RUNNER_TOKEN",
                deps = { "config" },
                env = {
                    value = "{{ .Last.Data.runner_token }}",
                },
            },
        },
    },
})

infra.dns({
    labels = { tf = "setup", dns = "1" },
    dc1 = true,
    global = false,
})
