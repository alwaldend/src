local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.vault_auth({ name = "default", approle = { name = "src_infra_download" } })

for _, stage in ipairs({ "local", "yandex", "dns" }) do
    lib.plugin_call({
        name = "tf_backend_" .. stage,
        plugin = "tf_backend",
        labels = { tf = stage },
        data = {
            vault_secret = "alwaldend.com/vault1/approles/src_infra_download/tf_backend/"
                .. stage,
            vault_secret_mount = "secrets",
        },
    })
end

infra.xo_login({
    labels = { xoa_login = "1" },
    oidc_provider = "src_infra_xcp_ng_provider",
})

infra.yc_auth({
    path = "yandex.cloud/org1/folders/src-infra-download/account_iam_key",
    labels = { tf = "yandex" },
})

infra.dns({ labels = { tf = "dns" }, dc1 = true })

infra.ansible_keys({
    labels = { ansible = "1" },
    vault_ssh = { backend = "ssh/clients/sign/admins", ttl = 60 * 60 * 2 },
})
