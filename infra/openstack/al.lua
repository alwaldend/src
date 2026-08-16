local infra = require("infra.al_lib")
local lib = require("al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_openstack",
    },
})

infra.ansible_keys({
    name = "ansible_keys",
    labels = { ansible = "1" },
    vault_ssh = {
        backend = "ssh/clients/sign/admins",
        ttl = 60 * 60 * 6,
    },
})

lib.plugin_call({
    name = "tf_backend_tf",
    plugin = "tf_backend",
    labels = { tf = "main" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_openstack/tf_backend/tf",
        vault_secret_mount = "secrets",
    },
})

infra.mikrotik({
    path = "alwaldend.com/vault1/approles/src_infra_openstack/mikrotik",
    host = "https://router1.dc1.alwaldend.com",
    labels = { tf = "main" },
})

lib.plugin_call({
    name = "kolla_passwords",
    plugin = "injector",
    labels = { kolla = "1" },
    data = {
        res = {
            {
                name = "kolla",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_openstack/kolla",
                    mount = "secrets",
                },
            },
            {
                name = "kolla_passwords_file",
                deps = { "kolla" },
                file = {
                    from_file = "infra/openstack/kolla/passwords.yml.tpl",
                },
            },
            {
                name = "KOLLA_PASSWORDS_FILE",
                deps = { "kolla_passwords_file" },
                env = {
                    value = "{{ index .Last.Files 0 }}",
                },
            },
        },
    },
})
