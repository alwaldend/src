local lib = require("al_lib")

lib.vault_auth({
    name = "xcp_ng",
    approle = { name = "src_infra_xcp_ng" },
})

lib.plugin_call({
    name = "tf_backend",
    plugin = "tf_backend",
    labels = { xcp_ng_tf = "1" },
    data = {
        vault_auth = "xcp_ng",
        vault_secret = "alwaldend.com/vault1/approles/src_infra_xcp_ng/tf_backend",
        vault_secret_mount = "secrets",
    },
})

lib.plugin_call({
    name = "xcpng_host_login",
    plugin = "injector",
    labels = { xcpng_host = "1" },
    data = {
        res = {
            {
                name = "xcpng_host",
                vault_auth = "xcp_ng",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_xcp_ng/host1",
                    mount = "secrets",
                },
            },
            {
                name = "XCPNG_USERNAME",
                deps = { "xcpng_host" },
                env = { value = "{{ .Last.Data.xcpng_username }}" },
            },
            {
                name = "XCPNG_PASSWORD",
                deps = { "xcpng_host" },
                env = { value = "{{ .Last.Data.xcpng_password }}" },
            },
        },
    },
})

lib.plugin_call({
    name = "xoa_ssh",
    plugin = "injector",
    labels = { xcp_ng_ansible = "1" },
    data = {
        res = {
            {
                name = "xoa_ssh",
                vault_auth = "xcp_ng",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_xcp_ng/xoa",
                    mount = "secrets",
                },
            },
            {
                name = "XOA_SSH_PASSWORD",
                deps = { "xoa_ssh" },
                env = { value = "{{ .Last.Data.xoa_ssh_password }}" },
            },
        },
    },
})

lib.plugin_call({
    name = "xoa_login",
    plugin = "injector",
    labels = { xoa = "1" },
    data = {
        res = {
            {
                name = "xoa",
                vault_auth = "xcp_ng",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_xcp_ng/xoa",
                    mount = "secrets",
                },
            },
            {
                name = "XOA_TOKEN",
                deps = { "xoa" },
                env = { value = "{{ .Last.Data.xoa_token }}" },
            },
        },
    },
})

lib.plugin_call({
    name = "env_xcp_ng",
    plugin = "injector",
    labels = { vault_env = "xcp_ng" },
    data = {
        res = {
            {
                name = "env_xcp_ng",
                vault_env = { conn = "default", auth = "xcp_ng" },
            },
        },
    },
})
