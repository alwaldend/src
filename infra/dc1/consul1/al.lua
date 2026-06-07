local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_dc1_consul1",
    },
})

lib.plugin_call({
    name = "tf_backend",
    plugin = "tf_backend",
    labels = { tf_setup = "1" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_dc1_consul1/tf_backend/tf_setup",
        vault_secret_mount = "secrets"
    },
})

lib.plugin_call({
    name = "tf_backend_main",
    plugin = "tf_backend",
    labels = { tf_main = "1" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_dc1_consul1/tf_backend/tf_main",
        vault_secret_mount = "secrets"
    },
})

lib.plugin_call({
    name = "pve_login",
    plugin = "pve_login",
    labels = { tf_setup = "1" },
})

lib.plugin_call({
    name = "consul_gossip",
    plugin = "injector",
    labels = { ansible = "1" },
    data = {
        res = {
            {
                name = "gossip_key",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_dc1_consul1/gossip",
                    mount = "secrets",
                },
            },
            {
                name = "AL_CONSUL_GOSSIP_KEY",
                deps = { "gossip_key" },
                env = {
                    value = "{{ .Last.Data.key }}",
                },
            }
        }
    }
})

lib.plugin_call({
    name = "consul",
    plugin = "injector",
    labels = { tf_main = "1", consul = "1" },
    data = {
        res = {
            {
                name = "CONSUL_HTTP_ADDR",
                env = {
                    value = "https://host1.consul1.dc1.alwaldend.com:8501"
                }
            },
            {
                name = "consul_bootstrap",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_dc1_consul1/bootstrap",
                    mount = "secrets",
                },
            },
            {
                name = "CONSUL_HTTP_TOKEN",
                deps = { "consul_bootstrap" },
                env = {
                    value = "{{ .Last.Data.SecretID }}",
                },
            },
            {
                name = "consul_cert",
                op = {
                    method = "write",
                    path = "pki/ica_servers/issue/src_infra_dc1_consul1_pki_server",
                    data = {
                        common_name = "src-infra-dc1-consul1-tf.consul1.dc1.alwaldend.com",
                        ttl = "3600",
                    },
                },
            },
            {
                name = "consul_cert_cert",
                deps = { "consul_cert" },
                file = {
                    value = "{{ .Last.Data.certificate }}",
                }
            },
            {
                name = "CONSUL_CLIENT_CERT",
                deps = { "consul_cert_cert" },
                env = {
                    value = "{{ index .Last.Files 0 }}",
                },
            },
            {
                name = "consul_cert_key",
                deps = { "consul_cert" },
                file = {
                    value = "{{ .Last.Data.private_key }}",
                }
            },
            {
                name = "CONSUL_CLIENT_KEY",
                deps = { "consul_cert_key" },
                env = {
                    value = "{{ index .Last.Files 0 }}",
                },
            },
        }
    }
})

infra.yc_auth({
    path = "yandex.cloud/org1/folders/src-infra-dc1-consul1/account_iam_key",
    labels = { tf_main = "1" },
})

infra.yc_bucket_auth({
    path = "yandex.cloud/org1/folders/src-infra-dc1-consul1/account_static_key",
    labels = { tf_main = "1" },
})

infra.yc_account({
    path = "yandex.cloud/org1/folders/src-infra-dc1-consul1/account",
    labels = { tf_main = "1" },
})
