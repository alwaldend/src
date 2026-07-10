local lib = require("al_lib")
local infra = require("infra.al_lib")

lib.vault_auth({
    name = "default",
    approle = {
        name = "src_infra_harbor",
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
        vault_secret = "alwaldend.com/vault1/approles/src_infra_harbor/tf_backend/tf_setup",
        vault_secret_mount = "secrets"
    },
})

lib.plugin_call({
    name = "tf_backend_tf",
    plugin = "tf_backend",
    labels = { tf = "main" },
    data = {
        vault_secret = "alwaldend.com/vault1/approles/src_infra_harbor/tf_backend/tf",
        vault_secret_mount = "secrets"
    },
})

lib.plugin_call({
    name = "pve_login",
    plugin = "pve_login",
    labels = { tf = "setup" },
})

lib.plugin_call({
    name = "harbor_login",
    plugin = "harbor_login",
    labels = { tf = "main" },
})

infra.k3s_token({
    name = "k3s_token",
    path = "alwaldend.com/vault1/approles/src_infra_harbor/config",
    labels = { ansible = "1" },
})

infra.kubernetes_login({
    name = "kubernetes_login",
    labels = { k8s = "1" },
    cluster_ca = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJkakNDQVIyZ0F3SUJBZ0lCQURBS0JnZ3Foa2pPUFFRREFqQWpNU0V3SHdZRFZRUUREQmhyTTNNdGMyVnkKZG1WeUxXTmhRREUzT0RNM01EWXhNVGt3SGhjTk1qWXdOekV3TVRZMU5URTVXaGNOTXpZd056QTNNVFkxTlRFNQpXakFqTVNFd0h3WURWUVFEREJock0zTXRjMlZ5ZG1WeUxXTmhRREUzT0RNM01EWXhNVGt3V1RBVEJnY3Foa2pPClBRSUJCZ2dxaGtqT1BRTUJCd05DQUFSQnlDb2VTM2diS3NiVTIxR2tFOWxta0RMcDlUWm1wYVR0Y1JidlhrRXMKNFBzVU5wTFZjS0RKWWwzOXp2U3NsSWNKM3JqQUUvdmRhZUdOc09EeXNja01vMEl3UURBT0JnTlZIUThCQWY4RQpCQU1DQXFRd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVWVWaTd1MVhPZE1TV0FsNHRsdTQzCmtsaklIcVF3Q2dZSUtvWkl6ajBFQXdJRFJ3QXdSQUlnYzB3RWhtVmRVeVU1MUlhTjRUNElNTENsbDJnbVF4eW0KaU0xSU5KcEVkUklDSUFJbi9xZzdoNW4zTUFLR3NGRm1tSXJ3UnVkWVkrTmRLd2VGS1hudXk3aDQKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
    oidc = {
        name = "src_infra_harbor_cluster_provider";
        scope = "openid user email groups";
        client_id = "d8LmH5d7CHKwOQ925CGllVDkytJ1mYgn";
        redirect_uri = "https://unused";
    },
})
