local infra = require("infra.al_lib")
local lib = require("al_lib")

infra.xo_login({
    labels = { xoa_login = "1" },
    oidc_provider = "src_infra_xcp_ng_provider",
})

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_openhands",
    },
})

infra.ansible_keys({
    name = "ansible_keys",
    labels = { ansible = "1" },
    vault_ssh = {
        backend = "ssh/clients/sign/admins",
        ttl = 60 * 60 * 2,
    },
})

-- Only tf_setup exists for this component, so one state backend covers it and
-- stores the lock and state as Vault KV entries under the AppRole path.
lib.plugin_call({
    name = "tf_backend",
    plugin = "tf_backend",
    labels = { tf = "setup", dns = "1" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_openhands/tf_backend",
        vault_secret_mount = "secrets",
    },
})

-- The agent server session key is the single browser-facing credential: the
-- canvas presents it to the agent server, and the automation server accepts
-- and presents the same value. Keeping one secret keeps those three agreeing
-- without duplicating the value across roles.
lib.plugin_call({
    name = "agent_server",
    plugin = "injector",
    labels = { ansible = "1" },
    data = {
        res = {
            {
                name = "agent_server",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_openhands/agent_server",
                    mount = "secrets",
                },
            },
            {
                name = "OPENHANDS_SESSION_API_KEY",
                deps = { "agent_server" },
                env = {
                    value = "{{ .Last.Data.session_api_key }}",
                },
            },
            {
                name = "OPENHANDS_SECRET_KEY",
                deps = { "agent_server" },
                env = {
                    value = "{{ .Last.Data.secret_key }}",
                },
            },
        },
    },
})

lib.plugin_call({
    name = "automation",
    plugin = "injector",
    labels = { ansible = "1" },
    data = {
        res = {
            {
                name = "automation",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_openhands/automation",
                    mount = "secrets",
                },
            },
            {
                name = "OPENHANDS_AUTOMATION_KV_SECRET",
                deps = { "automation" },
                env = {
                    value = "{{ .Last.Data.kv_secret }}",
                },
            },
        },
    },
})

infra.dns({
    labels = { tf = "setup", dns = "1" },
    dc1 = true,
})
