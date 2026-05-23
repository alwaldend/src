local lib = require("al_lib")

local M = {}

function M.tf_backend(t)
    local path, labels, name = t.path, t.labels, t.name or "tf_backend"
    lib.secret({
        name = name,
        labels = labels,
        kv = {
            path = path,
            mount = "secrets",
        },
    })
    lib.file({
        name = "backend_config",
        labels = labels,
        secrets = {name},
        value = [[
            bucket = "{{ .Secret.bucket }}"
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
    })
    lib.env({
        name = "AL_TF_BACKEND_CONFIG_1",
        labels = labels,
        files = {"backend_config"},
        value = "{{ .File.Path }}"
    })
end

function M.yc_auth(t)
    local path, labels, name, filename = t.path, t.labels, t.name or "yc_auth", t.filename or "service_account_key"
    lib.secret({
        name = name,
        labels = labels,
        kv = {
            path = path,
            mount = "secrets",
        },
    })
    lib.file({
        name = filename,
        labels = labels,
        secrets = {name},
        value = "{{ .Secret.service_account_key }}"
    })
    lib.env({
        name = "YC_CLOUD_ID",
        labels = labels,
        secrets = {name},
        value = "{{ .Secret.cloud_id }}",
    })
    lib.env({
        name = "TF_VAR_cloud_id",
        labels = labels,
        secrets = {name},
        value = "{{ .Secret.cloud_id }}",
    })
    lib.env({
        name = "YC_FOLDER_ID",
        labels = labels,
        secrets = {name},
        value = "{{ .Secret.folder_id }}",
    })
    lib.env({
        name = "TF_VAR_folder_id",
        labels = labels,
        secrets = {name},
        value = "{{ .Secret.folder_id }}",
    })
    lib.env({
        name = "YC_SERVICE_ACCOUNT_KEY_FILE",
        labels = labels,
        files = {filename},
        value = "{{ .File.Path }}",
    })
end

function M.yc_bucket_auth(t)
    local path, labels, name = t.path, t.labels, t.name or "yc_bucket_auth"
    lib.secret({
        name = name,
        labels = labels,
        kv = {
            path = path,
            mount = "secrets",
        },
    })
    lib.env({
        name = "AWS_ACCESS_KEY_ID",
        labels = labels,
        secrets = {name},
        value = "{{ .Secret.access_key }}",
    })
    lib.env({
        name = "AWS_SECRET_ACCESS_KEY",
        labels = labels,
        secrets = {name},
        value = "{{ .Secret.secret_key }}"
    })
end

function M.yc_account(t)
    local path, labels, name = t.path, t.labels, t.name or "yc_account"
    lib.secret({
        name = name,
        labels = labels,
        kv = {
            path = path,
            mount = "secrets",
        },
    })
    lib.env({
        name = "TF_VAR_service_account_id",
        labels = labels,
        secrets = {name},
        value = "{{ .Secret.service_account_id }}",
    })
end


function M.rclone_config(t)
    local name, labels, path = t.name or "rclone_config", t.labels, t.path
    lib.file({
        name = name,
        labels = labels,
        value = [[
            [remote]
            type = s3
            provider = AWS
            env_auth = true
            region = ru-central1
            endpoint = storage.yandexcloud.net
        ]]
    })
    lib.secret({
        name = "rclone_bucket",
        labels = labels,
        kv = {
            path = path,
            mount = "secrets",
        }
    })
    lib.env({
        name = "RCLONE_CONFIG",
        labels = labels,
        files = {name},
        value = "{{ .File.Path }}",
    })
    lib.env({
        name = "RCLONE_S3_BUCKET",
        labels = labels,
        secrets = {"rclone_bucket"},
        value = "{{ .Secret.bucket }}"
    })
end

return M
