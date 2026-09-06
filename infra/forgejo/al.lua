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
        name = "src_infra_dc1_forgejo1",
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

lib.plugin_call({
    name = "tf_backend_tf_setup",
    plugin = "tf_backend",
    labels = { tf = "setup" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_dc1_forgejo1/tf_backend/tf_setup",
        vault_secret_mount = "secrets",
    },
})

lib.plugin_call({
    name = "tf_backend_tf",
    plugin = "tf_backend",
    labels = { tf = "main" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_dc1_forgejo1/tf_backend/tf",
        vault_secret_mount = "secrets",
    },
})

lib.plugin_call({
    name = "forgejo_login",
    plugin = "forgejo_login",
    labels = { tf = "main" },
})

lib.plugin_call({
    name = "config",
    plugin = "injector",
    labels = { ansible = "1" },
    data = {
        res = {
            {
                name = "forgejo_oidc_client",
                op = {
                    method = "read",
                    path = "identity/oidc/client/src_infra_dc1_forgejo1_provider",
                },
            },
            {
                name = "FORGEJO_OIDC_CLIENT_ID",
                deps = { "forgejo_oidc_client" },
                env = { value = "{{ .Last.Data.client_id }}" },
            },
            {
                name = "FORGEJO_OIDC_CLIENT_SECRET",
                deps = { "forgejo_oidc_client" },
                env = { value = "{{ .Last.Data.client_secret }}" },
            },
            {
                name = "config",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_dc1_forgejo1/config",
                    mount = "secrets",
                },
            },
            {
                name = "LFS_JWT_SECRET",
                deps = { "config" },
                env = {
                    value = "{{ .Last.Data.lfs_jwt_secret }}",
                },
            },
            {
                name = "SECURITY_INTERNAL_TOKEN",
                deps = { "config" },
                env = {
                    value = "{{ .Last.Data.security_internal_token }}",
                },
            },
            {
                name = "OAUTH2_JWT_SECRET",
                deps = { "config" },
                env = {
                    value = "{{ .Last.Data.oauth2_jwt_secret }}",
                },
            },
        },
    },
})
