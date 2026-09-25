local lib = require("al_lib")

local M = {
    routeros_hosturl = "https://router1.dc1.alwaldend.com",
}

-- Register XO authentication with the caller's Vault OIDC provider.
-- Callers package //infra/xcp_ng/cmd/xo_login and select the same labels.
function M.xo_login(t)
    assert(
        type(t.oidc_provider) == "string" and t.oidc_provider ~= "",
        "xo_login requires oidc_provider"
    )
    lib.plugin({
        name = t.name or "xo_login",
        bin = "com_alwaldend_src/infra/xcp_ng/cmd/xo_login/xo_login_/xo_login",
        labels = t.labels,
        data = {
            xoa_url = "https://xoa.xcp-ng.alwaldend.com",
            discovery_url = "https://vault.alwaldend.com:8200/v1/identity/oidc/provider/"
                .. t.oidc_provider
                .. "/.well-known/openid-configuration",
        },
    })
end

-- Supply the configured endpoint for provider validation without logging in.
function M.pve_provider_inputs(t)
    lib.plugin_call({
        name = "pve_provider_inputs",
        plugin = "injector",
        labels = t.labels,
        data = {
            res = {
                {
                    name = "PM_API_URL",
                    env = { value = lib.pve_base_url .. "/api2/json" },
                },
            },
        },
    })
end

function M.k3s_token(t)
    local name, labels, path, mount =
        t.name or "k3s_token", t.labels, t.path, t.mount or "secrets"
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = {
            res = {
                {
                    name = name,
                    kv = {
                        path = path,
                        mount = mount,
                    },
                },
                {
                    name = "K3S_TOKEN",
                    deps = { name },
                    env = {
                        value = "{{ .Last.Data.k3s_token }}",
                    },
                },
            },
        },
    })
end

function M.kubernetes_login(t)
    local name, oidc, labels = t.name or "kubernetes_login", t.oidc, t.labels
    local cluster_ca = t.cluster_ca
    local res = {
        {
            name = name,
            oidc = oidc,
        },
        {
            name = name .. "_file",
            deps = { name },
            file = {
                extra = {
                    cluster_ca = cluster_ca,
                },
                value = [[
apiVersion: v1
clusters:
  - cluster:
      certificate-authority-data: "{{ .Extra.cluster_ca }}"
      server: https://127.0.0.1:6443
    name: default
contexts:
  - context:
      cluster: default
      user: default
    name: default
current-context: default
kind: Config
users:
  - name: default
    user:
      token: "{{ .Last.Data.id_token }}"
]],
            },
        },
        {
            name = "KUBECONFIG",
            deps = { name .. "_file" },
            env = {
                value = "{{ index .Last.Files 0 }}",
            },
        },
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res },
    })
end

function M.ansible_keys(t)
    local name, vault_ssh, labels =
        t.name or "ansible_keys", t.vault_ssh, t.labels
    local res = {
        {
            name = name,
            vault_ssh = vault_ssh,
        },
        {
            name = "ANSIBLE_PRIVATE_KEY_FILE",
            deps = { name },
            env = {
                value = "{{ .Last.Data.private_key }}",
            },
        },
        {
            name = "SSH_AUTH_SOCK",
            deps = { name },
            env = {
                value = "",
            },
        },
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res },
    })
end

function M.server_cert(t)
    local name, role, labels, data = t.name, t.role, t.labels, t.data
    local cert_name = "AL_SERVER_CERT_" .. name .. "_CERT"
    local key_name = "AL_SERVER_CERT_" .. name .. "_KEY"
    local res = {
        {
            name = name,
            op = {
                method = "write",
                data = data,
                path = "pki/ica_servers/issue/" .. role,
            },
        },
        {
            name = cert_name,
            deps = { name },
            file = {
                value = "{{ .VaultOp.certificate }}",
            },
        },
        {
            name = key_name,
            deps = { name },
            file = {
                value = "{{ .Last.Data.private_key }}",
            },
        },
        {
            name = cert_name,
            deps = { cert_name },
            env = {
                value = "{{ .Last.Data.Files[0] }}",
            },
        },
        {
            name = key_name,
            deps = { key_name },
            env = {
                value = "{{ .Last.Data.Files[0] }}",
            },
        },
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res },
    })
end

function M.acme_eab(t)
    local name, role, labels = t.name, t.role, t.labels
    local id_name = "AL_ACME_EAB_" .. name .. "_ID"
    local key_name = "AL_ACME_EAB_" .. name .. "_KEY"
    local res = {
        {
            name = name,
            op = {
                method = "write",
                path = "pki/ica_servers/roles/" .. role .. "/acme/new-eab",
            },
        },
        {
            name = id_name,
            deps = { name },
            env = {
                value = "{{ .VaultOp.id }}",
            },
        },
        {
            name = key_name,
            deps = { name },
            env = {
                value = "{{ .VaultOp.key }}",
            },
        },
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res },
    })
end

function M.tf_backend(t)
    local path, labels, name = t.path, t.labels, t.name or "tf_backend"
    local conf_name = name .. "_backend_config"
    local res = {
        {
            name = name,
            kv = {
                path = path,
                mount = "secrets",
            },
        },
        {
            name = conf_name,
            deps = { name },
            file = {
                value = [[
                    bucket = "{{ .Last.Data.bucket }}"
                    endpoints = {
                      s3 = "https://storage.yandexcloud.net"
                    }
                    region = "ru-central1"
                    use_lockfile = true
                    skip_region_validation      = true
                    skip_credentials_validation = true
                    skip_requesting_account_id  = true
                    skip_s3_checksum            = true
                ]],
            },
        },
        {
            name = "AL_TF_BACKEND_CONFIG_1",
            deps = { conf_name },
            env = {
                value = "{{ index .Last.Files 0 }}",
            },
        },
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res },
    })
end

function M.yc_auth(t)
    local path, labels, name, filename =
        t.path,
        t.labels,
        t.name or "yc_auth",
        t.filename or "service_account_key"
    local res = {
        {
            name = name,
            kv = {
                path = path,
                mount = "secrets",
            },
        },
        {
            name = filename,
            deps = { name },
            file = {
                value = "{{ .Last.Data.service_account_key }}",
            },
        },
        {
            name = "YC_CLOUD_ID",
            deps = { name },
            env = {
                value = "{{ .Last.Data.cloud_id }}",
            },
        },
        {
            name = "TF_VAR_cloud_id",
            deps = { name },
            env = {
                value = "{{ .Last.Data.cloud_id }}",
            },
        },
        {
            name = "YC_FOLDER_ID",
            deps = { name },
            env = {
                value = "{{ .Last.Data.folder_id }}",
            },
        },
        {
            name = "TF_VAR_folder_id",
            deps = { name },
            env = {
                value = "{{ .Last.Data.folder_id }}",
            },
        },
        {
            name = "YC_SERVICE_ACCOUNT_KEY_FILE",
            deps = { filename },
            env = {
                value = "{{ index .Last.Files 0 }}",
            },
        },
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res },
    })
end

function M.yc_bucket_auth(t)
    local path, labels, name = t.path, t.labels, t.name or "yc_bucket_auth"
    local res = {
        {
            name = name,
            kv = {
                path = path,
                mount = "secrets",
            },
        },
        {
            name = "AWS_ACCESS_KEY_ID",
            deps = { name },
            env = {
                value = "{{ .Last.Data.access_key }}",
            },
        },
        {
            name = "AWS_SECRET_ACCESS_KEY",
            deps = { name },
            env = {
                value = "{{ .Last.Data.secret_key }}",
            },
        },
    }
    lib.plugin_call({
        name = name,
        labels = labels,
        plugin = "injector",
        data = { res = res },
    })
end

function M.yc_account(t)
    local path, labels, name = t.path, t.labels, t.name or "yc_account"
    local res = {
        {
            name = name,
            kv = {
                path = path,
                mount = "secrets",
            },
        },
        {
            name = "TF_VAR_service_account_id",
            deps = { name },
            env = {
                value = "{{ .Last.Data.service_account_id }}",
            },
        },
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res },
    })
end

function M.rclone_config(t)
    local name, labels, path = t.name or "rclone_config", t.labels, t.path
    local res = {
        {
            name = name,
            file = {
                value = [[
                    [remote]
                    type = s3
                    provider = AWS
                    env_auth = true
                    region = ru-central1
                    endpoint = storage.yandexcloud.net
                ]],
            },
        },
        {
            name = "rclone_bucket",
            kv = {
                path = path,
                mount = "secrets",
            },
        },
        {
            name = "RCLONE_CONFIG",
            deps = { name },
            env = {
                value = "{{ index .Last.Files 0 }}",
            },
        },
        {
            name = "RCLONE_S3_BUCKET",
            deps = { "rclone_bucket" },
            env = {
                value = "{{ .Last.Data.bucket }}",
            },
        },
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res },
    })
end

function M.mikrotik(t)
    local name, labels, path, host, mount =
        t.name or "mikrotik", t.labels, t.path, t.host, t.mount or "secrets"
    local res = {
        {
            name = name,
            kv = {
                path = path,
                mount = mount,
            },
        },
        {
            name = "MIKROTIK_HOST",
            deps = { name },
            env = {
                value = host,
            },
        },
        {
            name = "MIKROTIK_USER",
            deps = { name },
            env = {
                value = "{{ .Last.Data.mikrotik_username }}",
            },
        },
        {
            name = "MIKROTIK_PASSWORD",
            deps = { name },
            env = {
                value = "{{ .Last.Data.mikrotik_password }}",
            },
        },
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res },
    })
end

-- Provider credentials for component-owned Terraform DNS records. Callers
-- select only the views they own; secret values remain in Vault.
function M.dns(t)
    local res = {}
    if t.global ~= false then
        res[#res + 1] = {
            name = "dns_cloudflare",
            vault_auth = t.vault_auth,
            kv = { path = "cloudflare.com/dns_token", mount = "secrets" },
        }
        res[#res + 1] = {
            name = "CLOUDFLARE_API_TOKEN",
            deps = { "dns_cloudflare" },
            env = { value = "{{ .Last.Data.cloudflare_api_token }}" },
        }
        res[#res + 1] = {
            name = "TF_VAR_dns_cloudflare_zone_id",
            deps = { "dns_cloudflare" },
            env = {
                value = '{{ with index .Last.Data "cloudflare_zone_id" }}{{ . }}{{ end }}',
            },
        }
    end
    if t.dc1 then
        res[#res + 1] = {
            name = "dns_mikrotik",
            vault_auth = t.vault_auth,
            kv = {
                path = "alwaldend.com/vault1/approles/src_infra_dns/mikrotik",
                mount = "secrets",
            },
        }
        res[#res + 1] = {
            name = "TF_VAR_dns_routeros_hosturl",
            deps = { "dns_mikrotik" },
            env = { value = M.routeros_hosturl },
        }
        res[#res + 1] = {
            name = "TF_VAR_dns_routeros_username",
            deps = { "dns_mikrotik" },
            env = { value = "{{ .Last.Data.username }}" },
        }
        res[#res + 1] = {
            name = "TF_VAR_dns_routeros_password",
            deps = { "dns_mikrotik" },
            env = { value = "{{ .Last.Data.password }}" },
        }
    end
    lib.plugin_call({
        name = t.name or "dns_providers",
        plugin = "injector",
        labels = t.labels,
        data = { res = res },
    })
end

return M
