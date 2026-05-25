local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.auth({
    name = "default",
    approle = {
        name = "src_infra_dc1_consul1",
    },
})

lib.secret({
    name = "gossip_key",
    kv = {
        path = "alwaldend.com/vault1/approles/src_infra_dc1_consul1/gossip",
        mount = "secrets",
    },
})

lib.secret({
    name = "bootstrap",
    kv = {
        path = "alwaldend.com/vault1/approles/src_infra_dc1_consul1/bootstrap",
        mount = "secrets",
    },
})

lib.env({
    name = "CONSUL_HTTP_TOKEN",
    secrets = { "bootstrap" },
    value = "{{ .Secret.SecretID }}",
    labels = { tf = "1" }
})

lib.vault_op({
    name = "consul_cert",
    method = "write",
    data = {
        common_name = "src-infra-dc1-consul1.tf.consul1.dc1.alwaldend.com",
        ttl = "3600",
    },
    path = "pki/ica_servers/issue/src_infra_dc1_consul1_pki_server",
})

lib.file({
    name = "consul_cert_cert",
    vault_ops = { "consul_cert" },
    value = "{{ .VaultOp.certificate }}",
})

lib.env({
    name = "CONSUL_CLIENT_CERT",
    files = { "consul_cert_cert" },
    value = "{{ .File.Path }}",
    labels = { tf = "1" }
})

lib.file({
    name = "consul_cert_key",
    vault_ops = { "consul_cert" },
    value = "{{ .VaultOp.private_key }}",
})

lib.env({
    name = "CONSUL_CLIENT_KEY",
    files = { "consul_cert_key" },
    value = "{{ .File.Path }}",
    labels = { tf = "1" }
})

lib.env({
    name = "AL_CONSUL_GOSSIP_KEY",
    secrets = { "gossip_key" },
    value = "{{ .Secret.key }}",
    labels = { ansible = "1" }
})

infra.tf_backend({
    path = "alwaldend.com/vault1/approles/src_infra_dc1_consul1/bucket",
    labels = { tf = "1" },
})

infra.yc_auth({
    path = "yandex.cloud/org1/folders/src-infra-dc1-consul1/account_iam_key",
    labels = { tf = "1" },
})

infra.yc_bucket_auth({
    path = "yandex.cloud/org1/folders/src-infra-dc1-consul1/account_static_key",
    labels = { tf = "1" },
})

infra.yc_account({
    path = "yandex.cloud/org1/folders/src-infra-dc1-consul1/account",
    labels = { tf = "1" },
})
