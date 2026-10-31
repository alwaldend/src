local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = { name = "user_simeonwarren" },
})

local ssh_signing = { backend = "ssh/clients/sign/admins", ttl = 60 * 60 * 2 }

-- Operator maintenance selects only SSH credentials and host trust.
local maintenance_server_ca_file = assert(io.open("data/ssh/ca_servers.pub"))
local maintenance_server_ca =
    assert(maintenance_server_ca_file:read("*a")):gsub("%s+$", "")
assert(maintenance_server_ca_file:close())

lib.plugin_call({
    name = "maintenance_ssh",
    plugin = "injector",
    labels = { maintenance = "1" },
    data = {
        res = {
            { name = "maintenance_key", vault_ssh = ssh_signing },
            {
                name = "maintenance_known_hosts",
                file = {
                    value = "@cert-authority *.alwaldend.com "
                        .. maintenance_server_ca
                        .. "\n",
                },
            },
            {
                name = "maintenance_config",
                deps = { "maintenance_key", "maintenance_known_hosts" },
                file = {
                    value = [[
Host *
    User ansible
    BatchMode yes
    ConnectTimeout 10
    ServerAliveInterval 15
    ServerAliveCountMax 2
    IdentityAgent none
    IdentitiesOnly yes
    IdentityFile "{{ .Res.maintenance_key.Data.private_key }}"
    CertificateFile "{{ .Res.maintenance_key.Data.public_key_signed }}"
    UserKnownHostsFile "{{ index .Res.maintenance_known_hosts.Files 0 }}"
    GlobalKnownHostsFile /dev/null
    StrictHostKeyChecking yes
    UpdateHostKeys no
    ControlMaster no
    ControlPath none
]],
                },
            },
            {
                name = "HOST_MAINTENANCE_SSH_KNOWN_HOSTS",
                deps = { "maintenance_known_hosts" },
                env = { value = "{{ index .Last.Files 0 }}" },
            },
            {
                name = "HOST_MAINTENANCE_SSH_CONFIG",
                deps = { "maintenance_config" },
                env = { value = "{{ index .Last.Files 0 }}" },
            },
            {
                name = "ANSIBLE_SSH_ARGS",
                deps = { "maintenance_config" },
                env = { value = '-F "{{ index .Last.Files 0 }}"' },
            },
            {
                name = "ANSIBLE_PRIVATE_KEY_FILE",
                deps = { "maintenance_key" },
                env = { value = "{{ .Last.Data.private_key }}" },
            },
        },
    },
})
