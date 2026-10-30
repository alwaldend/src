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

local ssh_signing = { backend = "ssh/clients/sign/admins", ttl = 60 * 60 * 2 }

-- Build-time input: the public host CA snapshot is owned by data/ssh.
local server_ca_file = assert(io.open("data/ssh/ca_servers.pub"))
local server_ca = assert(server_ca_file:read("*a")):gsub("%s+$", "")
assert(server_ca_file:close())

lib.plugin_call({
    name = "ssh",
    plugin = "injector",
    labels = { ssh = "1" },
    data = {
        res = {
            { name = "ssh_key", vault_ssh = ssh_signing },
            {
                name = "DOWNLOAD_SSH_PRIVATE_KEY",
                deps = { "ssh_key" },
                env = { value = "{{ .Last.Data.private_key }}" },
            },
            {
                name = "DOWNLOAD_SSH_CERTIFICATE",
                deps = { "ssh_key" },
                env = { value = "{{ .Last.Data.public_key_signed }}" },
            },
            {
                name = "ssh_known_hosts",
                file = {
                    value = "@cert-authority *.download.alwaldend.com "
                        .. server_ca
                        .. "\n",
                },
            },
            {
                name = "DOWNLOAD_SSH_KNOWN_HOSTS",
                deps = { "ssh_known_hosts" },
                env = { value = "{{ index .Last.Files 0 }}" },
            },
        },
    },
})
