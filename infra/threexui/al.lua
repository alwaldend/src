local lib = require("al_lib")
local infra = require("infra.al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_threexui",
    },
})

infra.ansible_keys({
    name = "ansible_keys",
    labels = { ansible = "1" },
    vault_ssh = {
        backend = "ssh/clients/sign/admins",
        ttl = 60 * 60 * 2
    }
})

lib.plugin_call({
    name = "tf_backend_tf_setup",
    plugin = "tf_backend",
    labels = { tf = "setup" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_threexui/tf_backend/tf_setup",
        vault_secret_mount = "secrets"
    },
})

lib.plugin_call({
    name = "tf_backend_tf",
    plugin = "tf_backend",
    labels = { tf = "main" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_threexui/tf_backend/tf",
        vault_secret_mount = "secrets"
    },
})

lib.plugin_call({
    name = "pve_login",
    plugin = "pve_login",
    labels = { tf = "setup" },
})

lib.plugin_call({
    name = "control_plane",
    plugin = "injector",
    labels = { ansible = "1", tf = "main" },
    data = {
        res = {
            {
                name = "control_plane",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_threexui/control_panel",
                    mount = "secrets"
                }
            },
            {
                name = "XUI_USERNAME",
                deps = { "control_plane" },
                env = {
                    value = "{{ .Last.Data.xui_username }}",
                }
            },
            {
                name = "XUI_PASSWORD",
                deps = { "control_plane" },
                env = {
                    value = "{{ .Last.Data.xui_password }}",
                }
            },
            {
                name = "TF_VAR_xui_url",
                deps = { "control_plane" },
                env = {
                    value = "{{ .Last.Data.xui_url }}",
                }
            },
            {
                name = "TF_VAR_xui_username",
                deps = { "control_plane" },
                env = {
                    value = "{{ .Last.Data.xui_username }}",
                }
            },
            {
                name = "TF_VAR_xui_password",
                deps = { "control_plane" },
                env = {
                    value = "{{ .Last.Data.xui_password }}",
                }
            },
            {
                name = "TF_VAR_xui_token",
                deps = { "njalla_node" },
                env = {
                    value = "{{ .Last.Data.xui_token }}",
                }
            },
        }
    }
})

lib.plugin_call({
    name = "njalla_node",
    plugin = "injector",
    labels = { ansible = "1", tf = "main" },
    data = {
        res = {
            {
                name = "njalla_node",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_threexui/node_njalla",
                    mount = "secrets"
                }
            },
            {
                name = "NJALLA_XUI_USERNAME",
                deps = { "njalla_node" },
                env = {
                    value = "{{ .Last.Data.xui_username }}",
                }
            },
            {
                name = "NJALLA_XUI_PASSWORD",
                deps = { "njalla_node" },
                env = {
                    value = "{{ .Last.Data.xui_password }}",
                }
            },
            {
                name = "NJALLA_XUI_BASE_PATH",
                deps = { "njalla_node" },
                env = {
                    value = "{{ .Last.Data.xui_base_path }}",
                }
            },
            {
                name = "TF_VAR_njalla_xui_url",
                deps = { "njalla_node" },
                env = {
                    value = "{{ .Last.Data.xui_url }}",
                }
            },
            {
                name = "TF_VAR_njalla_xui_username",
                deps = { "njalla_node" },
                env = {
                    value = "{{ .Last.Data.xui_username }}",
                }
            },
            {
                name = "TF_VAR_njalla_xui_password",
                deps = { "njalla_node" },
                env = {
                    value = "{{ .Last.Data.xui_password }}",
                }
            },
            {
                name = "TF_VAR_njalla_xui_token",
                deps = { "njalla_node" },
                env = {
                    value = "{{ .Last.Data.xui_token }}",
                }
            },
        }
    }
})
