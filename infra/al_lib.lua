local lib = require("al_lib")

local M = {}

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
                    cluster_ca = cluster_ca
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
]]
            },
        },
        {
            name = "KUBECONFIG",
            deps = { name .. "_file" },
            env = {
                value = "{{ index .Last.Files 0 }}",
            }
        },
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res }
    })
end

function M.ansible_keys(t)
    local name, vault_ssh, labels = t.name or "ansible_keys", t.vault_ssh, t.labels
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
            }
        },
        {
            name = "SSH_AUTH_SOCK",
            deps = { name },
            env = {
                value = "",
            }
        },
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res }
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
                path = "pki/ica_servers/issue/" .. role
            }
        },
        {
            name = cert_name,
            deps = { name },
            file = {
                value = "{{ .VaultOp.certificate }}"
            },
        },
        {
            name = key_name,
            deps = { name },
            file = {
                value = "{{ .Last.Data.private_key }}"
            },
        },
        {
            name = cert_name,
            deps = { cert_name },
            env = {
                value = "{{ .Last.Data.Files[0] }}"
            },
        },
        {
            name = key_name,
            deps = { key_name },
            env = {
                value = "{{ .Last.Data.Files[0] }}"
            }
        },
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res }
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
                path = "pki/ica_servers/roles/" .. role .. "/acme/new-eab"
            }
        },
        {
            name = id_name,
            deps = { name },
            env = {
                value = "{{ .VaultOp.id }}"
            }
        },
        {
            name = key_name,
            deps = {name},
            env = {
                value = "{{ .VaultOp.key }}"
            }
        }
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res }
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
            deps = {name},
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
                ]]
            }
        },
        {
            name = "AL_TF_BACKEND_CONFIG_1",
            deps = {conf_name},
            env  = {
                value = "{{ index .Last.Files 0 }}"
            }
        }
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res }
    })
end

function M.yc_auth(t)
    local path, labels, name, filename = t.path, t.labels, t.name or "yc_auth", t.filename or "service_account_key"
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
            deps = {name},
            file = {
                value = "{{ .Last.Data.service_account_key }}"
            }
        },
        {
            name = "YC_CLOUD_ID",
            deps = {name},
            env = {
                value = "{{ .Last.Data.cloud_id }}",
            }
        },
        {
            name = "TF_VAR_cloud_id",
            deps = {name},
            env = {
                value = "{{ .Last.Data.cloud_id }}",
            }
        },
        {
            name = "YC_FOLDER_ID",
            deps = {name},
            env = {
                value = "{{ .Last.Data.folder_id }}",
            }
        },
        {
            name = "TF_VAR_folder_id",
            deps = {name},
            env = {
                value = "{{ .Last.Data.folder_id }}",
            }
        },
        {
            name = "YC_SERVICE_ACCOUNT_KEY_FILE",
            deps = {filename},
            env = {
                value = "{{ index .Last.Files 0 }}",
            }
        },
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res }
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
            deps = {name},
            env = {
                value = "{{ .Last.Data.access_key }}",
            }
        },
        {
            name = "AWS_SECRET_ACCESS_KEY",
            deps = {name},
            env = {
                value = "{{ .Last.Data.secret_key }}"
            }
        },
    }
    lib.plugin_call({
        name = name,
        labels = labels,
        plugin = "injector",
        data = { res = res }
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
            deps = {name},
            env = {
                value = "{{ .Last.Data.service_account_id }}",
            }
        },
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res }
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
                ]]
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
            deps = {name},
            env = {
                value = "{{ index .Last.Files 0 }}",
            }
        },
        {
            name = "RCLONE_S3_BUCKET",
            deps = {"rclone_bucket"},
            env = {
                value = "{{ .Last.Data.bucket }}"
            }
        }
    }
    lib.plugin_call({
        name = name,
        plugin = "injector",
        labels = labels,
        data = { res = res }
    })
end

return M
