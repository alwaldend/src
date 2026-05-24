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
