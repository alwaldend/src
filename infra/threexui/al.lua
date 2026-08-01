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
    name = "pve_login",
    plugin = "pve_login",
    labels = { tf = "setup" },
})

infra.yc_auth({
    path = "yandex.cloud/org1/folders/src-infra-threexui/account_iam_key",
    labels = { tf = "setup" },
})

infra.yc_account({
    path = "yandex.cloud/org1/folders/src-infra-threexui/account",
    labels = { tf = "setup" },
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
    name = "mullvad_keys",
    plugin = "injector",
    labels = { tf = "main" },
    data = {
        res = {
            {
                name = "mullvad_keys",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_threexui/mullvad_keys",
                    mount = "secrets"
                }
            },
            {
                name = "TF_VAR_mullvad_keys",
                deps = { "mullvad_keys" },
                env = {
                    value = '{{ to_json_indent .Last.Data.keys "" "    " }}',
                }
            },
        },
    },
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
        }
    }
})

lib.plugin_call({
    name = "nodes",
    plugin = "injector",
    labels = { ansible = "1", tf = "main", sub = "1" },
    data = {
        res = {
            {
                name = "nodes",
                kv = {
                    path = "alwaldend.com/vault1/approles/src_infra_threexui/nodes",
                    mount = "secrets"
                }
            },
            {
                name = "NODE_XUI_USERNAME",
                deps = { "nodes" },
                env = {
                    value = "{{ .Last.Data.xui_username }}",
                }
            },
            {
                name = "TF_VAR_node_xui_username",
                deps = { "nodes" },
                env = {
                    value = "{{ .Last.Data.xui_username }}",
                }
            },
            {
                name = "NODE_XUI_PASSWORD",
                deps = { "nodes" },
                env = {
                    value = "{{ .Last.Data.xui_password }}",
                }
            },
            {
                name = "TF_VAR_node_xui_password",
                deps = { "nodes" },
                env = {
                    value = "{{ .Last.Data.xui_password }}",
                }
            },
            {
                name = "NODE_XUI_BASE_PATH",
                deps = { "nodes" },
                env = {
                    value = "{{ .Last.Data.xui_base_path }}",
                }
            },
            {
                name = "TF_VAR_node_xui_base_path",
                deps = { "nodes" },
                env = {
                    value = "{{ .Last.Data.xui_base_path }}",
                }
            },
            {
                name = "TF_VAR_node_xui_tokens",
                deps = { "nodes" },
                env = {
                    value = "{{ .Last.Data.xui_tokens | to_json }}",
                }
            },
        }
    }
})
